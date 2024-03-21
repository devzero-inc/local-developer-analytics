package collector

import (
	"context"
	"fmt"
	"lda/logging"
	"net"
	"os"
	"strings"
	"sync"

	"time"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
)

const SocketPath = "/tmp/lda.socket"

var (
	// Map to store ongoing commands
	ongoingCommands = make(map[string]Command)
	// Mutex to protect access to counter and conditionally starting/stopping collection
	collectionMutex sync.Mutex
	// Counter to track active commands
	activeCommandsCounter int
	// Context and cancel function to control the collection goroutine
	collectionContext    context.Context
	collectionCancelFunc context.CancelFunc = nil
	// Indicate if the collection is currently running
	isCollectionRunning bool = false
)

// Collect starts the collection of command and system information
func Collect() {
	logging.Log.Info().Msg("Collecting command and system information")

	// Create a context that listens for the interrupt signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Start collectSystemInformation in its own goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		collectSystemInformation(ctx, 120*time.Second, 0)
	}()

	// Start collectCommandInformation in its own goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := collectCommandInformation(); err != nil {
			logging.Log.Error().Err(err).Msg("Failed to collect command information")
			cancel()
		}
	}()

	// Wait for both functions to complete
	wg.Wait()

	logging.Log.Info().Msg("Collection stopped")
}

func collectSystemInformation(ctx context.Context, initialDuration, increaseDuration time.Duration) {
	// Perform initial collection
	if err := collectOnce(); err != nil {
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
			if err := collectOnce(); err != nil {
				logging.Log.Error().Err(err).Msg("Failed to collect system information")
			}
			// Increase the duration for the next tick
			currentDuration += increaseDuration
			logging.Log.Debug().Msgf("Next collection in %s", currentDuration)
		}
	}
}

func collectOnce() error {

	logging.Log.Debug().Msg("Collecting process")

	hostInfo, _ := host.Info()

	processes, err := process.Processes()
	if err != nil {
		logging.Log.Err(err).Msg("Error retrieving processes")
		return err
	}

	var processInfos []Process
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
			PID:            int(p.Pid),
			Name:           name,
			Status:         status,
			CreatedTime:    createTime,
			StoredTime:     time.Now().UnixMilli(),
			OS:             hostInfo.OS,
			Platform:       hostInfo.Platform,
			PlatformFamily: hostInfo.PlatformFamily,
			CPUUsage:       cpuPercent,
			MemoryUsage:    memorypercent,
		}

		InsertProcess(processInfo)
		processInfos = append(processInfos, processInfo)
	}

	return nil
}

func onStartCommand() {
	collectionMutex.Lock()
	defer collectionMutex.Unlock()

	activeCommandsCounter++
	// If the collection is not running, start it with a timeout
	if !isCollectionRunning {
		logging.Log.Debug().Msg("Starting collection")
		var timeoutDuration = 10 * time.Minute
		collectionContext, collectionCancelFunc = context.WithTimeout(context.Background(), timeoutDuration)
		go collectSystemInformation(collectionContext, 1*time.Second, 5)
		isCollectionRunning = true
	}
}

func onEndCommand() {
	collectionMutex.Lock()
	defer collectionMutex.Unlock()

	activeCommandsCounter--
	// If there are no more active commands, stop the collection
	if activeCommandsCounter == 0 && isCollectionRunning {
		logging.Log.Debug().Msg("Stopping collection")
		collectionCancelFunc()
		isCollectionRunning = false
	}
}

func collectCommandInformation() error {
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
			if err := handleSocketCollection(conn); err != nil {
				logging.Log.Error().Err(err).Msg("Error handling socket collection")
			}
		}()
	}
}

func handleSocketCollection(c net.Conn) error {
	defer c.Close()
	var buf [1024]byte
	n, err := c.Read(buf[:])
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
		if err := handleStartCommand(parts); err != nil {
			logging.Log.Error().Err(err).Msg("Error handling start command")
		}
	} else if parts[0] == "end" {
		if err := handleEndCommand(parts); err != nil {
			logging.Log.Error().Err(err).Msg("Error handling end command")
		}
	} else {
		logging.Log.Error().Msg("Invalid command format")
		return err
	}

	return nil
}

func handleStartCommand(parts []string) error {
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

	ongoingCommands[parts[4]] = command

	onStartCommand()

	return nil
}

func handleEndCommand(parts []string) error {

	if !IsCommandAcceptable(parts[1]) {
		logging.Log.Debug().Msg("Command is not acceptable")
		return fmt.Errorf("command is not acceptable")
	}

	logging.Log.Debug().Msgf("Parsing command: %s", parts[0])

	if command, exists := ongoingCommands[parts[4]]; exists {
		command.EndTime = time.Now().UnixMilli()
		command.ExecutionTime = command.EndTime - command.StartTime

		if err := InsertCommand(command); err != nil {
			logging.Log.Error().Err(err).Msg("Failed to insert command")
			return err
		}
		delete(ongoingCommands, parts[4])
		onEndCommand()
	} else {
		logging.Log.Error().Msg("Matching start command not found")
		return fmt.Errorf("matching start command not found")
	}

	return nil
}
