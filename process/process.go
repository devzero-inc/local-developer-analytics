package process

import (
	"errors"
	"lda/database"
	gen "lda/gen/api/v1"

	"github.com/rs/zerolog"
)

const (
	PsutilType = "psutil"
	PsType     = "ps"
)

// SystemProcess interface for process collection
type SystemProcess interface {
	Collect() ([]Process, error)
}

// Factory implementation for proc providers
type Factory struct {
	logger zerolog.Logger
}

// NewFactory init Factory implementation for proc providers
func NewFactory(logger zerolog.Logger) *Factory {
	return &Factory{
		logger: logger,
	}
}

// Create creates a new system process collector from the factory
func (f *Factory) Create(pType string) (SystemProcess, error) {
	switch pType {
	case PsutilType:
		return NewPsutil(f.logger), nil
	case PsType:
		return NewPs(f.logger), nil
	default:
		return nil, errors.New("system process type not supported")
	}
}

// Process is the model for process
type Process struct {
	Id   int64  `json:"id" db:"id"`
	PID  int64  `json:"pid" db:"pid"`
	Name string `json:"name" db:"name"`
	// R: Running; S: Sleep; T: Stop; I: Idle; Z: Zombie; W: Wait; L: Lock;
	Status         string  `json:"status" db:"status"`
	CreatedTime    int64   `json:"created_time" db:"created_time"`
	StoredTime     int64   `json:"stored_time" db:"stored_time"`
	OS             string  `json:"os" db:"os"`
	Platform       string  `json:"platform" db:"platform"`
	PlatformFamily string  `json:"platform_family" db:"platform_family"`
	CPUUsage       float64 `json:"cpu_usage" db:"cpu_usage"`
	MemoryUsage    float64 `json:"memory_usage" db:"memory_usage"`
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
		return nil, err
	}

	return processes, nil
}

// GetTopProcessesAndMetrics fetches the top processes based on a criterion like average CPU usage,
// and then fetches detailed time-series data for each top process.
func GetTopProcessesAndMetrics(start int64, end int64) (map[int64][]Process, error) {
	var topProcesses []Process
	processMetricsMap := make(map[int64][]Process)

	// Step 1: Identify the top processes
	topProcessesQuery := `SELECT name, pid, AVG(cpu_usage) AS cpu_usage, AVG(memory_usage) AS memory_usage
		FROM processes 
		WHERE stored_time BETWEEN ? AND ?
		GROUP BY name, pid
		ORDER BY cpu_usage DESC, memory_usage DESC
		LIMIT 20;`

	err := database.DB.Select(&topProcesses, topProcessesQuery, start, end)
	if err != nil {
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
			continue
		}

		processMetricsMap[process.PID] = metrics
	}

	return processMetricsMap, nil
}

// InsertProcess inserts a process into the database
func InsertProcess(process Process) error {
	query := `INSERT INTO processes (pid, name, status, created_time, stored_time, os, platform, platform_family, cpu_usage, memory_usage)
	VALUES (:pid, :name, :status, :created_time, :stored_time, :os, :platform, :platform_family, :cpu_usage, :memory_usage)`

	_, err := database.DB.NamedExec(query, process)
	return err
}

func MapProcessToProto(process Process) *gen.Process {
	return &gen.Process{
		Id:             process.Id,
		Pid:            process.PID,
		Name:           process.Name,
		Status:         process.Status,
		CreatedTime:    process.CreatedTime,
		StoredTime:     process.StoredTime,
		Os:             process.OS,
		Platform:       process.Platform,
		PlatformFamily: process.PlatformFamily,
		CpuUsage:       process.CPUUsage,
		MemoryUsage:    process.MemoryUsage,
	}
}
