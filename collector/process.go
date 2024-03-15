package collector

import (
	"lda/database"
	"lda/logging"
)

type Process struct {
	Id   int    `json:"id" db:"id"`
	PID  int    `json:"pid" db:"pid"`
	Name string `json:"name" db:"name"`
	// R: Running; S: Sleep; T: Stop; I: Idle; Z: Zombie; W: Wait; L: Lock;
	Status         string  `json:"status" db:"status"`
	StartTime      string  `json:"startTime" db:"start_time"`
	EndTime        string  `json:"endTime" db:"end_time"`
	ExecutionTime  int64   `json:"executionTime" db:"execution_time"`
	OS             string  `json:"os" db:"os"`
	Platform       string  `json:"platform" db:"platform"`
	PlatformFamily string  `json:"platformFamily" db:"platform_family"`
	CPUUsage       float64 `json:"cpuUsage" db:"cpu_usage"`
	UsedMemory     float32 `json:"usedMemory" db:"used_memory"`
}

func GetAllProceses() []Process {
	var processes []Process
	if err := database.DB.Select(&processes, "SELECT * FROM processes"); err != nil {
		logging.Log.Err(err).Msg("Failed to get all processes")
	}

	return processes
}

func GetAllProcessesForPeriod(start, end string) []Process {
	var processes []Process

	query := `SELECT * FROM processes WHERE start_time >= ? AND end_time <= ? ORDER BY start_time ASC`

	err := database.DB.Select(&processes, query, start, end)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to get processes with start and end times")
	}

	return processes
}

func InsertProcess(process Process) {
	query := `INSERT INTO processes (pid, name, status, start_time, end_time, execution_time, os, platform, platform_family, cpu_usage, used_memory)
	VALUES (:pid, :name, :status, :start_time, :end_time, :execution_time, :os, :platform, :platform_family, :cpu_usage, :used_memory)`

	_, err := database.DB.NamedExec(query, process)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to insert process")
	}
}
