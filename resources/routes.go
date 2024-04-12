package resources

import (
	"embed"
	"encoding/json"
	"lda/collector"
	"lda/logging"
	"lda/process"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

// Embedding directory
//
//go:embed views/*
var templateFS embed.FS

func showError(w http.ResponseWriter) {
	tmpl, err := template.ParseFS(templateFS, "views/error.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	loc, _ := time.LoadLocation("Local")
	now := time.Now().In(loc)

	var startMillis, endMillis int64

	start := r.URL.Query().Get("start")
	if start == "" {
		// Default start time to the start of today (00:00:00)
		startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		startMillis = startTime.UnixMilli()
	} else {
		// Parse the incoming start time and convert it to Unix milliseconds
		if parsedTime, err := time.ParseInLocation("2006-01-02T15:04", start, loc); err == nil {
			startMillis = parsedTime.UnixMilli()
		}
	}

	end := r.URL.Query().Get("end")
	if end == "" {
		// Default end time to the current time
		endMillis = now.UnixMilli()
	} else {
		// Parse the incoming end time and convert it to Unix milliseconds
		if parsedTime, err := time.ParseInLocation("2006-01-02T15:04", end, loc); err == nil {
			endMillis = parsedTime.UnixMilli()
		}
	}

	eStart := time.Now()

	commands, err := collector.GetAllCommandsForPeriod(startMillis, endMillis)
	if err != nil {
		showError(w)
		return
	}
	elapsed := time.Since(eStart)
	logging.Log.Debug().Msgf("Time taken to retrieve commands: %s", elapsed)

	eStart = time.Now()
	processes, err := process.GetAllProcessesForPeriod(startMillis, endMillis)
	if err != nil {
		showError(w)
		return
	}

	elapsed = time.Since(eStart)
	logging.Log.Debug().Msgf("Time taken to retrieve processes: %s", elapsed)

	eStart = time.Now()
	timeProcesses, err := process.GetTopProcessesAndMetrics(startMillis, endMillis)
	if err != nil {
		showError(w)
		return
	}

	elapsed = time.Since(eStart)
	logging.Log.Debug().Msgf("Time taken to retrieve time processes: %s", elapsed)

	commandsJson, err := PrepareCommandCategoriesExecutionTimeChartData(commands)
	if err != nil {
		showError(w)
		return
	}
	processResourceJson, err := PrepareProcessesResourceUsageChartData(processes)
	if err != nil {
		showError(w)
		return
	}
	cpuResourceJson, err := PrepareCPUTimeSeriesChartData(timeProcesses)
	if err != nil {
		showError(w)
		return
	}
	memoryResourceJson, err := PrepareMemoryTimeSeriesChartData(timeProcesses)
	if err != nil {
		showError(w)
		return
	}

	tmpl, err := template.ParseFS(templateFS, "views/index.html")
	if err != nil {
		showError(w)
		return
	}

	logging.Log.Info().Msgf("Rendering template with start: %s, end: %s", start, end)

	if start == "" {
		start = time.UnixMilli(startMillis).UTC().Format("2006-01-02T15:04")
	}

	if end == "" {
		end = time.UnixMilli(endMillis).UTC().Format("2006-01-02T15:04")
	}

	if err := tmpl.Execute(w, map[string]interface{}{
		"CommandsJSON":         commandsJson,
		"ProcessesJSON":        processResourceJson,
		"CPUTimeSeriesJSON":    cpuResourceJson,
		"MemoryTimeSeriesJSON": memoryResourceJson,
		"StartTime":            start,
		"EndTime":              end,
	}); err != nil {
		showError(w)
	}
}

func commandHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	label := queryParams.Get("label")

	if label == "" {
		showError(w)
		return
	}

	loc, _ := time.LoadLocation("Local")
	now := time.Now().In(loc)

	var startMillis, endMillis int64

	start := queryParams.Get("start")
	if start == "" {
		// Default start time to the start of today (00:00:00)
		startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		startMillis = startTime.UnixMilli()
	} else {
		// Parse the incoming start time and convert it to Unix milliseconds
		if parsedTime, err := time.ParseInLocation("2006-01-02T15:04", start, loc); err == nil {
			startMillis = parsedTime.UnixMilli()
		}
	}

	end := queryParams.Get("end")
	if end == "" {
		// Default end time to the current time
		endMillis = now.UnixMilli()
	} else {
		// Parse the incoming end time and convert it to Unix milliseconds
		if parsedTime, err := time.ParseInLocation("2006-01-02T15:04", end, loc); err == nil {
			endMillis = parsedTime.UnixMilli()
		}
	}

	commands, err := collector.GetAllCommandsForCategoryForPeriod(
		label, startMillis, endMillis)
	if err != nil {
		showError(w)
		return
	}

	commandsJson, err := PrepareCommandsExecutionTimeChartData(commands)
	if err != nil {
		showError(w)
		return
	}

	tmpl, err := template.ParseFS(templateFS, "views/command.html")
	if err != nil {
		showError(w)
		return
	}

	logging.Log.Info().Msgf("Rendering template with start: %s, end: %s", start, end)

	if start == "" {
		start = time.UnixMilli(startMillis).UTC().Format("2006-01-02T15:04")
	}

	if end == "" {
		end = time.UnixMilli(endMillis).UTC().Format("2006-01-02T15:04")
	}

	if err := tmpl.Execute(w, map[string]interface{}{
		"CommandsJSON": commandsJson,
		"StartTime":    start,
		"EndTime":      end,
		"Commands":     commands,
	}); err != nil {
		showError(w)
	}
}

func overviewHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	label := queryParams.Get("id")

	if label == "" {
		showError(w)
		return
	}

	i, err := strconv.ParseInt(label, 10, 64)
	if err != nil {
		showError(w)
		return
	}

	command, err := collector.GetCommandById(i)
	if err != nil {
		showError(w)
		return
	}

	processes, err := process.GetAllProcessesForPeriod(command.StartTime, command.EndTime)
	if err != nil {
		showError(w)
		return
	}

	processesJson, err := json.Marshal(processes)
	if err != nil {
		showError(w)
		return
	}

	timeProcesses, err := process.GetTopProcessesAndMetrics(command.StartTime, command.EndTime)
	if err != nil {
		showError(w)
		return
	}

	processResourceJson, err := PrepareProcessesResourceUsageChartData(processes)
	if err != nil {
		showError(w)
		return
	}
	cpuResourceJson, err := PrepareCPUTimeSeriesChartData(timeProcesses)
	if err != nil {
		showError(w)
		return
	}
	memoryResourceJson, err := PrepareMemoryTimeSeriesChartData(timeProcesses)
	if err != nil {
		showError(w)
		return
	}

	tmpl, err := template.ParseFS(templateFS, "views/overview.html")
	if err != nil {
		showError(w)
		return
	}

	if err := tmpl.Execute(w, map[string]interface{}{
		"ProcessResourceJSON":  processResourceJson,
		"CPUTimeSeriesJSON":    cpuResourceJson,
		"MemoryTimeSeriesJSON": memoryResourceJson,
		"Processes":            processes,
		"ProcessJSON":          string(processesJson),
	}); err != nil {
		showError(w)
	}
}

// Serve registers the HTTP handlers for the application
func Serve() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/command", commandHandler)
	http.HandleFunc("/overview", overviewHandler)
}
