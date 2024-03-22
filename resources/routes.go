package resources

import (
	"embed"
	"lda/collector"
	"lda/logging"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

// Embedding directory
//
//go:embed views/*
var templateFS embed.FS

func Serve() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

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

		commands, err := collector.GetAllCommandsForPeriod(startMillis, endMillis)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		processes, err := collector.GetAllProcessesForPeriod(startMillis, endMillis)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		timeProcesses, err := collector.GetTopProcessesAndMetrics(startMillis, endMillis)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		commandsJson, err := PrepareCommandCategoriesExecutionTimeChartData(commands)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		processResourceJson, err := PrepareProcessesResourceUsageChartData(processes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cpuResourceJson, err := PrepareCPUTimeSeriesChartData(timeProcesses)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		memoryResourceJson, err := PrepareMemoryTimeSeriesChartData(timeProcesses)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFS(templateFS, "views/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/command", func(w http.ResponseWriter, r *http.Request) {

		queryParams := r.URL.Query()

		label := queryParams.Get("label")

		if label == "" {
			http.Error(w, "Pull parameter is required", http.StatusBadRequest)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		commandsJson, err := PrepareCommandsExecutionTimeChartData(commands)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFS(templateFS, "views/command.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	http.HandleFunc("/overview", func(w http.ResponseWriter, r *http.Request) {

		queryParams := r.URL.Query()

		label := queryParams.Get("id")

		if label == "" {
			http.Error(w, "Pull parameter is required", http.StatusBadRequest)
		}

		i, err := strconv.ParseInt(label, 10, 64)
		if err != nil {
			http.Error(w, "Failed to parse id", http.StatusBadRequest)
		}

		command, err := collector.GetCommandById(i)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		processes, err := collector.GetAllProcessesForPeriod(command.StartTime, command.EndTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		timeProcesses, err := collector.GetTopProcessesAndMetrics(command.StartTime, command.EndTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		processResourceJson, err := PrepareProcessesResourceUsageChartData(processes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cpuResourceJson, err := PrepareCPUTimeSeriesChartData(timeProcesses)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		memoryResourceJson, err := PrepareMemoryTimeSeriesChartData(timeProcesses)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFS(templateFS, "views/overview.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, map[string]interface{}{
			"ProcessesJSON":        processResourceJson,
			"CPUTimeSeriesJSON":    cpuResourceJson,
			"MemoryTimeSeriesJSON": memoryResourceJson,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})
}
