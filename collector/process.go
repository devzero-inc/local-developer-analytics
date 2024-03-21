package collector

import (
	"lda/database"
	"lda/logging"
)

// Process is the model for process
type Process struct {
	Id   int    `json:"id" db:"id"`
	PID  int    `json:"pid" db:"pid"`
	Name string `json:"name" db:"name"`
	// R: Running; S: Sleep; T: Stop; I: Idle; Z: Zombie; W: Wait; L: Lock;
	Status         string  `json:"status" db:"status"`
	CreatedTime    int64   `json:"created_time" db:"created_time"`
	StoredTime     int64   `json:"stored_time" db:"stored_time"`
	OS             string  `json:"os" db:"os"`
	Platform       string  `json:"platform" db:"platform"`
	PlatformFamily string  `json:"platform_family" db:"platform_family"`
	CPUUsage       float64 `json:"cpu_usage" db:"cpu_usage"`
	MemoryUsage    float32 `json:"memory_usage" db:"memory_usage"`
}

// GetAllProcesses fetches all processes from the database
func GetAllProcesses() ([]Process, error) {
	var processes []Process
	if err := database.DB.Select(&processes, "SELECT * FROM processes"); err != nil {
		logging.Log.Err(err).Msg("Failed to get all processes")
		return nil, err
	}

	return processes, nil
}

// GetAllProcessesForPeriod fetches all processes for a given period
func GetAllProcessesForPeriod(start int64, end int64) ([]Process, error) {
	var processes []Process

	query := `SELECT pid, name, AVG(cpu_usage) AS cpu_usage, AVG(memory_usage) AS memory_usage 
              FROM processes
              WHERE stored_time >= ? AND stored_time <= ?
              GROUP BY pid, name
              ORDER BY AVG(cpu_usage) DESC, AVG(memory_usage) DESC
              LIMIT 500`

	err := database.DB.Select(&processes, query, start, end)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to get aggregated processes with start and end times")
		return nil, err
	}

	return processes, nil
}

// GetTopProcessesAndMetrics fetches the top processes based on a criterion like average CPU usage,
// and then fetches detailed time-series data for each top process.
func GetTopProcessesAndMetrics(start int64, end int64) (map[int][]Process, error) {
	var topProcesses []Process
	processMetricsMap := make(map[int][]Process)

	// Step 1: Identify the top processes
	topProcessesQuery := `SELECT name, pid, AVG(cpu_usage) AS cpu_usage, AVG(memory_usage) AS memory_usage
		FROM processes 
		WHERE stored_time BETWEEN ? AND ?
		GROUP BY name, pid
		ORDER BY cpu_usage DESC, memory_usage DESC
		LIMIT 20;`

	err := database.DB.Select(&topProcesses, topProcessesQuery, start, end)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to get top processes")
		return nil, err
	}

	// Step 2: Fetch time-series data for each top process
	for _, process := range topProcesses {
		var metrics []Process

		metricsQuery := `SELECT name, pid, cpu_usage,  memory_usage, stored_time
			FROM processes
			WHERE name = ? AND pid = ? AND stored_time BETWEEN ? AND ?
			ORDER BY stored_time DESC
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

// InsertProcess inserts a process into the database
func InsertProcess(process Process) {
	query := `INSERT INTO processes (pid, name, status, created_time, stored_time, os, platform, platform_family, cpu_usage, memory_usage)
	VALUES (:pid, :name, :status, :created_time, :stored_time, :os, :platform, :platform_family, :cpu_usage, :memory_usage)`

	_, err := database.DB.NamedExec(query, process)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to insert process")
	}
}
