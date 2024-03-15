package cmd

import (
	"context"
	"lda/database"
	"lda/logging"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"time"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
)

type Process struct {
	Id   int    `json:"id" db:"id"`
	PID  int    `json:"pid" db:"pid"`
	Name string `json:"name" db:"name"`
	// R: Running; S: Sleep; T: Stop; I: Idle; Z: Zombie; W: Wait; L: Lock;
	Status         string  `json:"status" db:"status"`
	StartTime      string  `json:"startTime" db:"start_time"`
	ExecutionTime  string  `json:"executionTime" db:"execution_time"`
	OS             string  `json:"os" db:"os"`
	Platform       string  `json:"platform" db:"platform"`
	PlatformFamily string  `json:"platformFamily" db:"platform_family"`
	CPUUsage       float64 `json:"cpuUsage" db:"cpu_usage"`
	UsedMemory     float32 `json:"usedMemory" db:"used_memory"`
}

func GetAllProceses() {
	var processes []Process
	if err := database.DB.Select(&processes, "SELECT * FROM processes"); err != nil {
		logging.Log.Err(err).Msg("Failed to get all processes")
	}
}

func InsertProcess(process Process) {
	query := `INSERT INTO processes (pid, name, status, start_time, execution_time, os, platform, platform_family, cpu_usage, used_memory)
VALUES (:pid, :name, :status, :start_time, :execution_time, :os, :platform, :platform_family, :cpu_usage, :used_memory)`

	_, err := database.DB.NamedExec(query, process)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to insert process")
	}
}

type Command struct {
	Id            int    `json:"id" db:"id"`
	PID           int    `json:"pid" db:"pid"`
	Command       string `json:"command" db:"command"`
	User          string `json:"user" db:"user"`
	Directory     string `json:"directory" db:"directory"`
	ExecutionTime string `json:"executionTime" db:"execution_time"`
	StartTime     string `json:"startTime" db:"start_time"`
	EndTime       string `json:"endTime" db:"end_time"`
}

func GetAllCommands() {
	var commands []Command
	if err := database.DB.Select(&commands, "SELECT * FROM commands"); err != nil {
		logging.Log.Err(err).Msg("Failed to get all commands")
	}
}

func InsertCommand(command Command) {
	query := `INSERT INTO commands (command, user, directory, execution_time, start_time, end_time)
VALUES (:command, :user, :directory, :execution_time, :start_time, :end_time)`

	_, err := database.DB.NamedExec(query, command)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to insert command")
	}
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
	logging.Log.Info().Msg("Collecting command and system information")

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

				processInfo := Process{
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
				InsertProcess(processInfo)
				processInfos = append(processInfos, processInfo)
			}

			//jsonData, err := json.MarshalIndent(processInfos, "", "    ")
			//if err != nil {
			//	logging.Log.Err(err).Msg("Error marshalling data to JSON")
			//	continue
			//}

			//logging.Log.Info().Msg(string(jsonData))
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

			data := string(buf[:n])
			parts := strings.Split(data, "|")

			if len(parts) != 6 {
				logging.Log.Error().Msg("Invalid command format")
			} else {
				command := Command{
					Command:       parts[0],
					Directory:     parts[1],
					User:          parts[2],
					StartTime:     parts[3],
					EndTime:       parts[4],
					ExecutionTime: parts[5],
				}

				InsertCommand(command)

				logging.Log.Info().Msgf("Received: %s", string(buf[:n]))

			}

		}(conn)
	}
}
