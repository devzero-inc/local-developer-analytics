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
	StartTime      int64   `json:"startTime" db:"start_time"`
	EndTime        int64   `json:"endTime" db:"end_time"`
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

func GetAllProcessesForPeriod(start int64, end int64) []Process {
	var procesess []Process

	query := `SELECT pid, name, AVG(cpu_usage) AS cpu_usage, AVG(used_memory) AS used_memory 
              FROM processes
              WHERE end_time >= ? AND end_time <= ?
              GROUP BY pid, name
              ORDER BY AVG(cpu_usage) DESC, AVG(used_memory) DESC
              LIMIT 500`

	err := database.DB.Select(&procesess, query, start, end)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to get aggregated processes with start and end times")
	}

	return procesess
}

// GetTopProcessesAndMetrics fetches the top processes based on a criterion like average CPU usage,
// and then fetches detailed time-series data for each top process.
func GetTopProcessesAndMetrics(start int64, end int64) (map[int][]Process, error) {
	var topProcesses []Process
	processMetricsMap := make(map[int][]Process)

	// Step 1: Identify the top processes
	topProcessesQuery := `SELECT name, pid, AVG(cpu_usage) AS cpu_usage, AVG(used_memory) AS used_memory
		FROM processes 
		WHERE end_time BETWEEN ? AND ?
		GROUP BY name, pid
		ORDER BY cpu_usage DESC, used_memory DESC
		LIMIT 20;`

	err := database.DB.Select(&topProcesses, topProcessesQuery, start, end)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to get top processes")
		return nil, err
	}

	// Step 2: Fetch time-series data for each top process
	for _, process := range topProcesses {
		var metrics []Process

		metricsQuery := `SELECT name, pid, cpu_usage,  used_memory, end_time
			FROM processes
			WHERE name = ? AND pid = ? AND end_time BETWEEN ? AND ?
			ORDER BY end_time DESC
			LIMIT 20;`

		err := database.DB.Select(&metrics, metricsQuery, process.Name, process.PID, start, end)
		if err != nil {
			logging.Log.Err(err).Msgf("Failed to get time-series data for process %d", process.PID)
			continue
		}

		processMetricsMap[process.PID] = metrics
	}

	return processMetricsMap, nil
}

//func GetAllProcessesForPeriod(start int64, end int64) []Process {
//	var processes []Process
//
//	query := `SELECT * FROM processes WHERE start_time >= ? AND end_time <= ? ORDER BY start_time ASC`
//
//	err := database.DB.Select(&processes, query, start, end)
//	if err != nil {
//		logging.Log.Err(err).Msg("Failed to get processes with start and end times")
//	}
//
//	return processes
//}

func InsertProcess(process Process) {
	query := `INSERT INTO processes (pid, name, status, start_time, end_time, execution_time, os, platform, platform_family, cpu_usage, used_memory)
	VALUES (:pid, :name, :status, :start_time, :end_time, :execution_time, :os, :platform, :platform_family, :cpu_usage, :used_memory)`

	_, err := database.DB.NamedExec(query, process)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to insert process")
	}
}
