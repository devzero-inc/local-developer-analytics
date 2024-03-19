package collector

import (
	"context"
	"lda/logging"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"time"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
)

const SocketPath = "/tmp/lda.socket"

var (
	// Map to store ongoing commands
	ongoingCommands = make(map[int64]Command)
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
		collectSystemInformation(ctx, 120*time.Second)
	}()

	// Start collectCommandInformation in its own goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := collectCommandInformation(); err != nil {
			logging.Log.Error().Err(err).Msg("Failed to collect command information")
			cancel() // Optionally cancel the context if there's an error
		}
	}()

	// Wait for both functions to complete
	wg.Wait()

	logging.Log.Info().Msg("Collection stoped")
}

func collectSystemInformation(ctx context.Context, tickerDuration time.Duration) {
	logging.Log.Debug().Msgf("Collecting system information every %s", tickerDuration)

	// Perform initial collection before starting the ticker
	collectOnce()

	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop() // Ensure ticker is stopped to avoid leaks

	for {
		select {
		case <-ctx.Done():
			logging.Log.Info().Msg("Shutting down collection of system information")
			return
		case <-ticker.C:
			// Perform the collection on each tick
			collectOnce()
		}
	}
}

func collectOnce() {

	logging.Log.Debug().Msg("Collecting process")

	hostInfo, _ := host.Info()

	processes, err := process.Processes()
	if err != nil {
		logging.Log.Err(err).Msg("Error retrieving processes")
		return
	}

	var processInfos []Process
	for _, p := range processes {
		createTime, _ := p.CreateTime()
		now := time.Now()
		// Calculate executionTime directly in milliseconds as an int64
		executionTimeMs := int64(now.UnixNano()/1e6) - createTime

		name, _ := p.Name()

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
			StartTime:      createTime,
			EndTime:        now.UnixMilli(),
			ExecutionTime:  executionTimeMs,
			OS:             hostInfo.OS,
			Platform:       hostInfo.Platform,
			PlatformFamily: hostInfo.PlatformFamily,
			CPUUsage:       cpuPercent,
			UsedMemory:     memorypercent,
		}

		InsertProcess(processInfo)
		processInfos = append(processInfos, processInfo)
	}
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
		go collectSystemInformation(collectionContext, 1*time.Second)
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

	logging.Log.Info().Msg("Listening on " + SocketPath)
	for {
		conn, err := listener.Accept()
		if err != nil {
			logging.Log.Error().Err(err).Msg("Failed to accept connection")
			continue
		}

		// Handle each connection in a separate goroutine
		go func(c net.Conn) {
			defer c.Close()
			var buf [1024]byte
			n, err := c.Read(buf[:])
			if err != nil {
				logging.Log.Error().Err(err).Msg("Error reading from socket")
				return
			}

			data := string(buf[:n])
			parts := strings.Split(data, "|")

			logging.Log.Debug().Msgf("Received: %s", string(buf[:n]))

			if parts[0] == "start" {

				if len(parts) != 5 {
					logging.Log.Error().Msg("Invalid command format")
					return // Make sure to return after logging the error to prevent further execution
				}

				// Use strings.TrimSpace to remove any leading/trailing whitespace or newlines
				startTimeString := strings.TrimSpace(parts[4])

				startTimeTimestamp, err := strconv.ParseInt(startTimeString, 10, 64)
				if err != nil {
					logging.Log.Error().Err(err).Msg("Invalid start time format")
					return
				}

				logging.Log.Debug().Msgf("Parsing command: %s", parts[0])

				command := Command{
					Category:  ParseCommand(parts[1]),
					Command:   parts[1],
					Directory: parts[2],
					User:      parts[3],
					StartTime: startTimeTimestamp,
				}

				ongoingCommands[command.StartTime] = command

				// Call onStartCommand to increment the activeCommandsCounter and start the collection if necessary
				onStartCommand()

				//InsertCommand(command)
			} else if parts[0] == "end" {

				if len(parts) != 7 {
					logging.Log.Error().Msg("Invalid command format")
					return // Make sure to return after logging the error to prevent further execution
				}

				// Use strings.TrimSpace to remove any leading/trailing whitespace or newlines
				executionTimeString := strings.TrimSpace(parts[6])
				startTimeString := strings.TrimSpace(parts[4])
				endTimeString := strings.TrimSpace(parts[5])

				executionTimeMs, err := strconv.ParseInt(executionTimeString, 10, 64)
				if err != nil {
					logging.Log.Error().Err(err).Msg("Invalid execution time format")
					return
				}

				startTimeTimestamp, err := strconv.ParseInt(startTimeString, 10, 64)
				if err != nil {
					logging.Log.Error().Err(err).Msg("Invalid start time format")
					return
				}

				endTimeTimestamp, err := strconv.ParseInt(endTimeString, 10, 64)
				if err != nil {
					logging.Log.Error().Err(err).Msg("Invalid end time format")
					return
				}

				logging.Log.Debug().Msgf("Parsing command: %s", parts[0])

				if command, exists := ongoingCommands[startTimeTimestamp]; exists {
					command.EndTime = endTimeTimestamp
					command.ExecutionTime = executionTimeMs

					InsertCommand(command)

					delete(ongoingCommands, startTimeTimestamp)
					onEndCommand()
				} else {
					logging.Log.Error().Msg("Matching start command not found")
					return
				}

			} else {
				logging.Log.Error().Msg("Invalid command format")
				return // Make sure to return after logging the error to prevent further execution
			}
		}(conn)
	}
}
