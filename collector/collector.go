package collector

import (
	"context"
	"fmt"
	"lda/client"
	"lda/config"
	gen "lda/gen/api/v1"
	"lda/logging"
	"net"
	"os"
	"strings"
	"sync"

	"time"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
)

// TODO move this to /var/run or other appropriate location based on OS.
const SocketPath = "/tmp/lda.socket"

// Collector collects command and system information
type Collector struct {
	socketPath            string
	ongoingCommands       map[string]Command
	collectionMutex       sync.Mutex
	activeCommandsCounter int
	collectionContext     context.Context
	collectionCancelFunc  context.CancelFunc
	isCollectionRunning   bool
	client                *client.Client
}

// NewCollector creates a new collector instance
func NewCollector(socketPath string, client *client.Client) *Collector {
	return &Collector{
		socketPath:      socketPath,
		ongoingCommands: make(map[string]Command),
		client:          client,
	}
}

// Collect starts the collection of command and system information
func (c *Collector) Collect() {
	logging.Log.Info().Msg("Collecting command and system information")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.collectSystemInformation(ctx, time.Duration(config.AppConfig.ProcessInterval)*time.Second, 0)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.collectCommandInformation(); err != nil {
			logging.Log.Error().Err(err).Msg("Failed to collect command information")
			cancel()
		}
	}()

	wg.Wait()

	logging.Log.Info().Msg("Collection stopped")
}

func (c *Collector) collectSystemInformation(ctx context.Context, initialDuration time.Duration, increaseDuration time.Duration) {
	// Perform initial collection
	if err := c.collectOnce(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to collect system information")
	}

	currentDuration := initialDuration

	for {
		select {
		case <-ctx.Done():
			logging.Log.Info().Msg("Shutting down collection of system information")
			return
		case <-time.After(currentDuration):
			// Perform the collection on each tick
			if err := c.collectOnce(); err != nil {
				logging.Log.Error().Err(err).Msg("Failed to collect system information")
			}
			// Increase the duration for the next tick
			currentDuration += increaseDuration
			logging.Log.Debug().Msgf("Next collection in %s", currentDuration)
		}
	}
}

func (c *Collector) collectOnce() error {

	logging.Log.Debug().Msg("Collecting process")

	hostInfo, _ := host.Info()

	processes, err := process.Processes()
	if err != nil {
		logging.Log.Err(err).Msg("Error retrieving processes")
		return err
	}

	var processInfos []Process
	var processMetrics []*gen.Process
	for _, p := range processes {
		createTime, err := p.CreateTime()
		if err != nil {
			logging.Log.Err(err).Msg("Error retrieving create time")
			continue
		}

		name, err := p.Name()
		if err != nil {
			logging.Log.Err(err).Msg("Error retrieving name")
			continue
		}

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			logging.Log.Err(err).Msg("Error retrieving CPU percent")
			continue
		}

		memorypercent, err := p.MemoryPercent()
		if err != nil {
			logging.Log.Err(err).Msg("Error retrieving memory percent")
			continue
		}

		status, err := p.Status()
		if err != nil {
			logging.Log.Err(err).Msg("Error retrieving status")
			continue
		}

		// Adjust the Process struct to accept ExecutionTime as int64 if not already
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
		processInfos = append(processInfos, processInfo)

		if c.client != nil {
			processMetrics = append(processMetrics, MapProcessToProto(processInfo))
		}
	}

	if c.client != nil {
		if err := c.client.SendProcesses(processMetrics); err != nil {
			logging.Log.Error().Err(err).Msg("Failed to send processes")
		}
	}

	return nil
}

func (c *Collector) onStartCommand() {
	c.collectionMutex.Lock()
	defer c.collectionMutex.Unlock()

	c.activeCommandsCounter++
	// If the collection is not running, start it with a timeout
	if !c.isCollectionRunning {
		logging.Log.Debug().Msg("Starting collection")
		var timeoutDuration = 10 * time.Minute
		c.collectionContext, c.collectionCancelFunc = context.WithTimeout(context.Background(), timeoutDuration)
		go c.collectSystemInformation(
			c.collectionContext,
			time.Duration(config.AppConfig.CommandInterval)*time.Second,
			time.Duration(config.AppConfig.CommandIntervalMultiplier)*time.Second,
		)
		c.isCollectionRunning = true
	}
}

func (c *Collector) onEndCommand() {
	c.collectionMutex.Lock()
	defer c.collectionMutex.Unlock()

	c.activeCommandsCounter--
	// If there are no more active commands, stop the collection
	if c.activeCommandsCounter == 0 && c.isCollectionRunning {
		logging.Log.Debug().Msg("Stopping collection")
		c.collectionCancelFunc()
		c.isCollectionRunning = false
	}
}

func (c *Collector) collectCommandInformation() error {
	if err := os.RemoveAll(SocketPath); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to clean up existing socket")
		return err
	}

	listener, err := net.Listen("unix", SocketPath)
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to listen on UNIX socket")
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logging.Log.Error().Err(err).Msg("Failed to accept connection")
			continue
		}

		// Handle each connection in a separate goroutine
		go func() {
			if err := c.handleSocketCollection(conn); err != nil {
				logging.Log.Error().Err(err).Msg("Error handling socket collection")
			}
		}()
	}
}

func (c *Collector) handleSocketCollection(con net.Conn) error {
	defer con.Close()
	var buf [1024]byte
	n, err := con.Read(buf[:])
	if err != nil {
		logging.Log.Error().Err(err).Msg("Error reading from socket")
		return err
	}

	data := string(buf[:n])
	parts := strings.Split(data, "|")

	logging.Log.Debug().Msgf("Received: %s", string(buf[:n]))

	if len(parts) != 5 {
		logging.Log.Error().Msg("Invalid command format")
		return fmt.Errorf("invalid command format")
	}

	if parts[0] == "start" {
		if err := c.handleStartCommand(parts); err != nil {
			logging.Log.Error().Err(err).Msg("Error handling start command")
		}
	} else if parts[0] == "end" {
		if err := c.handleEndCommand(parts); err != nil {
			logging.Log.Error().Err(err).Msg("Error handling end command")
		}
	} else {
		logging.Log.Error().Msg("Invalid command format")
		return err
	}

	return nil
}

func (c *Collector) handleStartCommand(parts []string) error {
	if !IsCommandAcceptable(parts[1]) {
		logging.Log.Debug().Msg("Command is not acceptable")
		return fmt.Errorf("command is not acceptable")
	}

	logging.Log.Debug().Msgf("Parsing command: %s", parts[0])

	command := Command{
		Category:  ParseCommand(parts[1]),
		Command:   parts[1],
		Directory: parts[2],
		User:      parts[3],
		StartTime: time.Now().UnixMilli(),
	}

	c.ongoingCommands[parts[4]] = command

	c.onStartCommand()

	return nil
}

func (c *Collector) handleEndCommand(parts []string) error {

	if !IsCommandAcceptable(parts[1]) {
		logging.Log.Debug().Msg("Command is not acceptable")
		return fmt.Errorf("command is not acceptable")
	}

	logging.Log.Debug().Msgf("Parsing command: %s", parts[0])

	if command, exists := c.ongoingCommands[parts[4]]; exists {
		command.EndTime = time.Now().UnixMilli()
		command.ExecutionTime = command.EndTime - command.StartTime

		if err := InsertCommand(command); err != nil {
			logging.Log.Error().Err(err).Msg("Failed to insert command")
			return err
		}
		delete(c.ongoingCommands, parts[4])
		c.onEndCommand()

		if c.client != nil {
			if err := c.client.SendCommands([]*gen.Command{MapCommandToProto(command)}); err != nil {
				logging.Log.Error().Err(err).Msg("Failed to send command")
			}
		}
	} else {
		logging.Log.Error().Msg("Matching start command not found")
		return fmt.Errorf("matching start command not found")
	}

	return nil
}
