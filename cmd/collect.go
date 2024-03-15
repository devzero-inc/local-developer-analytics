package cmd

import (
	"context"
	"encoding/json"
	"lda/logging"
	"net"
	"os"
	"sync"

	"github.com/spf13/cobra"

	"time"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
)

type ProcessInfo struct {
	PID  int    `json:"pid"`
	Name string `json:"name"`
	// R: Running; S: Sleep; T: Stop; I: Idle; Z: Zombie; W: Wait; L: Lock;
	Status         string  `json:"status"`
	StartTime      string  `json:"startTime"`
	ExecutionTime  string  `json:"executionTime"`
	OS             string  `json:"os"`
	Platform       string  `json:"platform"`
	PlatformFamily string  `json:"platformFamily"`
	CPUUsage       float64 `json:"cpuUsage"`
	UsedMemory     float32 `json:"usedMemory"`
}

var (
	collectCmd = &cobra.Command{
		Use:   "collect",
		Short: "Collect command and system information",
		Long:  `Collect and process command and system information.`,

		Run: collect,
	}
)

func collect(_ *cobra.Command, _ []string) {
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

			var processInfos []ProcessInfo
			for _, p := range processes {
				createTime, _ := p.CreateTime()
				startTime := time.Unix(createTime/1000, 0)
				executionTime := time.Since(startTime).Round(time.Second)
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

				processInfo := ProcessInfo{
					PID:            int(p.Pid),
					Name:           name,
					Status:         status,
					StartTime:      startTime.Format("2006-01-02 15:04:05"),
					ExecutionTime:  executionTime.String(),
					OS:             hostInfo.OS,
					Platform:       hostInfo.Platform,
					PlatformFamily: hostInfo.PlatformFamily,
					CPUUsage:       cpuPercent,
					UsedMemory:     memorypercent,
				}
				processInfos = append(processInfos, processInfo)
			}

			jsonData, err := json.MarshalIndent(processInfos, "", "    ")
			if err != nil {
				logging.Log.Err(err).Msg("Error marshalling data to JSON")
				continue
			}

			logging.Log.Info().Msg(string(jsonData))
		}
	}
}

func collectCommandInformation() error {
	socketPath := "/tmp/myapp.socket"
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
			logging.Log.Info().Msgf("Received: %s", string(buf[:n]))
		}(conn)
	}
}
