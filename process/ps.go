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

type Ps struct {
	logger zerolog.Logger
}

func NewPs(logger zerolog.Logger) *Ps {
	return &Ps{
		logger: logger,
	}
}

func (p *Ps) Collect() ([]Process, error) {
	p.logger.Debug().Msg("Collecting process")

	cmd := exec.Command("ps", "axo", "pid,pcpu,pmem,lstart,comm")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(&out)
	scanner.Scan() // Skipping the header line

	var processInfo []Process
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		// Parsing the first three fields as PID, CPU and MEM
		pid, _ := strconv.ParseInt(fields[0], 10, 64)
		cpuUsage, _ := strconv.ParseFloat(fields[1], 64)
		memUsage, _ := strconv.ParseFloat(fields[2], 64)

		// Reconstruct lstart from fields 3 to 7
		lstart := strings.Join(fields[3:8], " ")
		const lstartLayout = "Mon Jan 2 15:04:05 2006"
		startTime, err := time.Parse(lstartLayout, lstart)
		if err != nil {
			p.logger.Err(err).Msg("Error parsing start time")
			continue
		}

		// The command name is the rest, starting from field 8
		// This assumes that the command name is the last field and can contain spaces
		name := strings.Join(fields[8:], " ")

		processInfo = append(processInfo, Process{
			PID:         pid,
			Name:        path.Base(name),
			CPUUsage:    cpuUsage,
			MemoryUsage: memUsage,
			CreatedTime: startTime.UnixMilli(),
			StoredTime:  time.Now().UnixMilli(),
			OS:          runtime.GOOS,
			Platform:    runtime.GOOS,
		})

		p.logger.Debug().Msgf("PID: %d, CPU: %f, MEM: %f, Start: %s, Name: %s", pid, cpuUsage, memUsage, startTime, name)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return processInfo, nil
}