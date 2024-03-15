package collector

import (
	"context"
	"encoding/json"
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

const socketPath = "/tmp/lda.socket"

func Collect() {

	// Create a context that listens for the interrupt signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Start collectSystemInformation in its own goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		collectSystemInformation(ctx)
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
}

func collectSystemInformation(ctx context.Context) {

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop() // Ensure ticker is stopped to avoid leaks

	for {
		select {
		case <-ctx.Done():
			logging.Log.Info().Msg("Shutting down collection of system information")
			return
		case <-ticker.C:
			hostInfo, _ := host.Info()

			processes, err := process.Processes()
			if err != nil {
				logging.Log.Err(err).Msg("Error retrieving processes")
				continue
			}

			var processInfos []Process
			for _, p := range processes {
				createTime, _ := p.CreateTime()
				// No need to convert to time.Time object if only calculating duration in ms
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
					StartTime:      time.Unix(createTime/1000, (createTime%1000)*1e6).Format("2006-01-02 15:04:05"),
					EndTime:        now.Format("2006-01-02 15:04:05"),
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

			jsonData, err := json.MarshalIndent(processInfos, "", "    ")
			if err != nil {
				logging.Log.Err(err).Msg("Error marshalling data to JSON")
				continue
			}

			logging.Log.Debug().Msg(string(jsonData))
		}
	}
}

func collectCommandInformation() error {
	if err := os.RemoveAll(socketPath); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to clean up existing socket")
		return err
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to listen on UNIX socket")
		return err
	}
	defer listener.Close()

	logging.Log.Info().Msg("Listening on " + socketPath)
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

			if len(parts) != 6 {
				logging.Log.Error().Msg("Invalid command format")
				return // Make sure to return after logging the error to prevent further execution
			}

			// Use strings.TrimSpace to remove any leading/trailing whitespace or newlines
			executionTimeString := strings.TrimSpace(parts[5])

			// Now convert the cleaned-up string to int64
			executionTimeMs, err := strconv.ParseInt(executionTimeString, 10, 64)
			if err != nil {
				logging.Log.Error().Err(err).Msg("Invalid execution time format")
				return
			}

			command := Command{
				Command:       parts[0],
				Directory:     parts[1],
				User:          parts[2],
				StartTime:     parts[3],
				EndTime:       parts[4],
				ExecutionTime: executionTimeMs, // Use the converted int64 value
			}

			InsertCommand(command)

		}(conn)
	}
}
