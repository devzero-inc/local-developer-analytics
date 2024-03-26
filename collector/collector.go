package collector

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"lda/client"
	gen "lda/gen/api/v1"
	"net"
	"os"
	"strings"
	"sync"

	"time"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
)

// TODO move this to /var/run or other appropriate location based on OS,
// TODO /var/run has issues with permisisons os have to explore a bit more.
const SocketPath = "/tmp/lda.socket"

// Collector collects command and system information
type Collector struct {
	socketPath       string
	client           *client.Client
	logger           zerolog.Logger
	collectionConfig collectionConfig
	intervalConfig   IntervalConfig
}

// IntervalConfig contains the configuration for the collection intervals
type IntervalConfig struct {
	ProcessInterval           time.Duration
	CommandInterval           time.Duration
	CommandIntervalMultiplier time.Duration
	MaxConcurrentCommands     int
}

// collectionConfig contains the configuration for the collection process
type collectionConfig struct {
	ongoingCommands       map[string]Command
	collectionMutex       sync.Mutex
	activeCommandsCounter int
	collectionContext     context.Context
	collectionCancelFunc  context.CancelFunc
	isCollectionRunning   bool
}

// NewCollector creates a new collector instance
func NewCollector(socketPath string, client *client.Client, logger zerolog.Logger, config IntervalConfig) *Collector {
	return &Collector{
		socketPath: socketPath,
		client:     client,
		logger:     logger,
		collectionConfig: collectionConfig{
			ongoingCommands: make(map[string]Command),
		},
		intervalConfig: config,
	}
}

// Collect starts the collection of command and system information
func (c *Collector) Collect() {
	c.logger.Info().Msg("Collecting command and system information")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.collectSystemInformation(ctx, c.intervalConfig.ProcessInterval*time.Second, 0)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.collectCommandInformation(); err != nil {
			c.logger.Error().Err(err).Msg("Failed to collect command information")
			cancel()
		}
	}()

	wg.Wait()

	c.logger.Info().Msg("Collection stopped")
}

func (c *Collector) collectSystemInformation(ctx context.Context, initialDuration time.Duration, increaseDuration time.Duration) {
	// Perform initial collection
	if err := c.collectOnce(); err != nil {
		c.logger.Error().Err(err).Msg("Failed to collect system information")
	}

	currentDuration := initialDuration

	for {
		select {
		case <-ctx.Done():
			c.logger.Debug().Msg("Shutting down collection of system information")
			return
		case <-time.After(currentDuration):
			// Perform the collection on each tick
			if err := c.collectOnce(); err != nil {
				c.logger.Error().Err(err).Msg("Failed to collect system information")
			}
			// Increase the duration for the next tick
			currentDuration += increaseDuration
			c.logger.Debug().Msgf("Next collection in %s", currentDuration)
		}
	}
}

func (c *Collector) collectOnce() error {

	c.logger.Debug().Msg("Collecting process")

	hostInfo, _ := host.Info()

	processes, err := process.Processes()
	if err != nil {
		c.logger.Err(err).Msg("Error retrieving processes")
		return err
	}

	var processMetrics []*gen.Process
	for _, p := range processes {
		createTime, err := p.CreateTime()
		if err != nil {
			c.logger.Err(err).Msg("Error retrieving create time")
			continue
		}

		name, err := p.Name()
		if err != nil {
			c.logger.Err(err).Msg("Error retrieving name")
			continue
		}

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			c.logger.Err(err).Msg("Error retrieving CPU percent")
			continue
		}

		memorypercent, err := p.MemoryPercent()
		if err != nil {
			c.logger.Err(err).Msg("Error retrieving memory percent")
			continue
		}

		status, err := p.Status()
		if err != nil {
			c.logger.Err(err).Msg("Error retrieving status")
			continue
		}

		processInfo := Process{
			PID:            int64(p.Pid),
			Name:           name,
			Status:         status,
			CreatedTime:    createTime,
			StoredTime:     time.Now().UnixMilli(),
			OS:             hostInfo.OS,
			Platform:       hostInfo.Platform,
			PlatformFamily: hostInfo.PlatformFamily,
			CPUUsage:       cpuPercent,
			MemoryUsage:    float64(memorypercent),
		}

		InsertProcess(processInfo)

		if c.client != nil {
			processMetrics = append(processMetrics, MapProcessToProto(processInfo))
		}
	}

	if c.client != nil {
		if err := c.client.SendProcesses(processMetrics); err != nil {
			c.logger.Error().Err(err).Msg("Failed to send processes")
		}
	}

	return nil
}

func (c *Collector) onStartCommand() {
	c.collectionConfig.collectionMutex.Lock()
	defer c.collectionConfig.collectionMutex.Unlock()

	c.collectionConfig.activeCommandsCounter++
	// If the collection is not running, start it with a timeout
	if !c.collectionConfig.isCollectionRunning {
		c.logger.Debug().Msg("Starting collection")
		var timeoutDuration = 10 * time.Minute
		c.collectionConfig.collectionContext, c.collectionConfig.collectionCancelFunc =
			context.WithTimeout(context.Background(), timeoutDuration)
		go c.collectSystemInformation(
			c.collectionConfig.collectionContext,
			c.intervalConfig.CommandInterval*time.Second,
			c.intervalConfig.CommandIntervalMultiplier*time.Second,
		)
		c.collectionConfig.isCollectionRunning = true
	}
}

func (c *Collector) onEndCommand() {
	c.collectionConfig.collectionMutex.Lock()
	defer c.collectionConfig.collectionMutex.Unlock()

	c.collectionConfig.activeCommandsCounter--
	// If there are no more active commands, stop the collection
	if c.collectionConfig.activeCommandsCounter == 0 && c.collectionConfig.isCollectionRunning {
		c.logger.Debug().Msg("Stopping collection")
		c.collectionConfig.collectionCancelFunc()
		c.collectionConfig.isCollectionRunning = false
	}
}

func (c *Collector) collectCommandInformation() error {
	if err := os.RemoveAll(SocketPath); err != nil {
		c.logger.Error().Err(err).Msg("Failed to clean up existing socket")
		return err
	}

	listener, err := net.Listen("unix", SocketPath)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to listen on UNIX socket")
		return err
	}
	defer listener.Close()

	// Limit the number of concurrent goroutines handling connections
	semaphore := make(chan struct{}, c.intervalConfig.MaxConcurrentCommands)

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-ctx.Done() // Wait for context cancellation
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				// If the context is canceled, stop accepting new connections
				return nil
			default:
				c.logger.Error().Err(err).Msg("Failed to accept connection")
				continue
			}
		}

		semaphore <- struct{}{} // Acquire
		go func(conn net.Conn) {
			defer func() {
				<-semaphore // Release
			}()
			if err := c.handleSocketCollection(conn); err != nil {
				c.logger.Error().Err(err).Msg("Error handling socket collection")
			}
		}(conn)
	}
}

func (c *Collector) handleSocketCollection(con net.Conn) error {
	defer con.Close()
	var buf [1024]byte
	n, err := con.Read(buf[:])
	if err != nil {
		c.logger.Error().Err(err).Msg("Error reading from socket")
		return err
	}

	data := string(buf[:n])
	parts := strings.Split(data, "|")

	c.logger.Debug().Msgf("Received: %s", string(buf[:n]))

	if len(parts) != 5 {
		c.logger.Error().Msg("Invalid command format")
		return fmt.Errorf("invalid command format")
	}

	if parts[0] == "start" {
		if err := c.handleStartCommand(parts); err != nil {
			c.logger.Error().Err(err).Msg("Error handling start command")
		}
	} else if parts[0] == "end" {
		if err := c.handleEndCommand(parts); err != nil {
			c.logger.Error().Err(err).Msg("Error handling end command")
		}
	} else {
		c.logger.Error().Msg("Invalid command format")
		return err
	}

	return nil
}

func (c *Collector) handleStartCommand(parts []string) error {
	if !IsCommandAcceptable(parts[1]) {
		c.logger.Debug().Msg("Command is not acceptable")
		return fmt.Errorf("command is not acceptable")
	}

	c.logger.Debug().Msgf("Parsing command: %s", parts[0])

	command := Command{
		Category:  ParseCommand(parts[1]),
		Command:   parts[1],
		Directory: parts[2],
		User:      parts[3],
		StartTime: time.Now().UnixMilli(), // TODO: there are some issues with sending time through shell because of ms support on MAC, explore more
	}

	c.collectionConfig.ongoingCommands[parts[4]] = command

	c.onStartCommand()

	return nil
}

func (c *Collector) handleEndCommand(parts []string) error {

	if !IsCommandAcceptable(parts[1]) {
		c.logger.Debug().Msg("Command is not acceptable")
		return fmt.Errorf("command is not acceptable")
	}

	c.logger.Debug().Msgf("Parsing command: %s", parts[0])

	if command, exists := c.collectionConfig.ongoingCommands[parts[4]]; exists {
		command.EndTime = time.Now().UnixMilli()
		command.ExecutionTime = command.EndTime - command.StartTime

		if err := InsertCommand(command); err != nil {
			c.logger.Error().Err(err).Msg("Failed to insert command")
			return err
		}
		delete(c.collectionConfig.ongoingCommands, parts[4])
		c.onEndCommand()

		if c.client != nil {
			if err := c.client.SendCommands([]*gen.Command{MapCommandToProto(command)}); err != nil {
				c.logger.Error().Err(err).Msg("Failed to send command")
			}
		}
	} else {
		c.logger.Error().Msg("Matching start command not found")
		return fmt.Errorf("matching start command not found")
	}

	return nil
}
