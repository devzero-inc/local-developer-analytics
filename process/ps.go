package process

import (
	"bufio"
	"bytes"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// Ps is the type for the ps process collector
type Ps struct {
	logger zerolog.Logger
}

// NewPs creates a new Ps instance
func NewPs(logger zerolog.Logger) *Ps {
	return &Ps{
		logger: logger,
	}
}

// Collect collects the process information using the ps command
func (p *Ps) Collect() ([]Process, error) {
	p.logger.Debug().Msg("Collecting process")

	// Adjust the command to include PPID
	cmd := exec.Command("ps", "axo", "pid,ppid,pcpu,pmem,lstart,comm")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(&out)
	scanner.Scan() // Skip the header line

	var processInfo []Process

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		// Parse PID, PPID, CPU and MEM usage
		pid, _ := strconv.ParseInt(fields[0], 10, 64)
		ppid, _ := strconv.ParseInt(fields[1], 10, 64)
		cpuUsage, _ := strconv.ParseFloat(fields[2], 64)
		memUsage, _ := strconv.ParseFloat(fields[3], 64)

		// Parse the start time
		lstart := strings.Join(fields[4:9], " ")
		const lstartLayout = "Mon Jan 2 15:04:05 2006"
		startTime, err := time.Parse(lstartLayout, lstart)
		if err != nil {
			p.logger.Err(err).Msg("Error parsing start time")
			continue
		}

		// Command name might contain spaces, so we join remaining fields
		name := strings.Join(fields[9:], " ")

		// Create the Process instance
		process := Process{
			PID:         pid,
			PPID:        ppid,
			Name:        path.Base(name),
			CPUUsage:    cpuUsage,
			MemoryUsage: memUsage,
			CreatedTime: startTime.UnixMilli(),
			StoredTime:  time.Now().UnixMilli(),
			OS:          runtime.GOOS,
			Platform:    runtime.GOOS,
		}

		// Append to the list of processes
		processInfo = append(processInfo, process)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return processInfo, nil
}
