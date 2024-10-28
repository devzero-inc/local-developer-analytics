package process

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
)

// Psutil is the type for the psutil process collector
type Psutil struct {
	logger zerolog.Logger
}

// NewPsutil creates a new Psutil instance
func NewPsutil(logger zerolog.Logger) *Psutil {
	return &Psutil{
		logger: logger,
	}
}

// Collect collects the process information using the psutil library
func (p *Psutil) Collect() ([]Process, error) {

	p.logger.Debug().Msg("Collecting process")

	hostInfo, _ := host.Info()

	processes, err := process.Processes()
	if err != nil {
		p.logger.Err(err).Msg("Error retrieving processes")
		return nil, err
	}

	var processInfo []Process
	for _, proc := range processes {
		createTime, err := proc.CreateTime()
		if err != nil {
			p.logger.Err(err).Msg("Error retrieving create time")
			continue
		}

		name, err := proc.Name()
		if err != nil {
			p.logger.Err(err).Msg("Error retrieving name")
			continue
		}

		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			p.logger.Err(err).Msg("Error retrieving CPU percent")
			continue
		}

		memorypercent, err := proc.MemoryPercent()
		if err != nil {
			p.logger.Err(err).Msg("Error retrieving memory percent")
			continue
		}

		status, err := proc.Status()
		if err != nil {
			p.logger.Err(err).Msg("Error retrieving status")
			continue
		}

		ppid, err := proc.Ppid()
		if err != nil {
			p.logger.Err(err).Msg("Error retrieving parent PID")
		}

		processInfo = append(processInfo, Process{
			PID:            int64(proc.Pid),
			PPID:           int64(ppid),
			Name:           name,
			Status:         status,
			CreatedTime:    createTime,
			StoredTime:     time.Now().UnixMilli(),
			OS:             hostInfo.OS,
			Platform:       hostInfo.Platform,
			PlatformFamily: hostInfo.PlatformFamily,
			CPUUsage:       cpuPercent,
			MemoryUsage:    float64(memorypercent),
		})
	}

	return processInfo, nil
}
