package resources

import (
	"embed"
	"encoding/json"
	"lda/collector"
	"lda/logging"
	"net/http"
	"text/template"
	"time"
)

// Embedding directory
//go:embed views/*
var templateFS embed.FS

func Serve() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		loc, _ := time.LoadLocation("Local")
		now := time.Now().In(loc)

		// Function to convert time.Time to Unix milliseconds
		toUnixMillis := func(t time.Time) int64 {
			return t.UnixNano() / int64(time.Millisecond)
		}

		var startMillis, endMillis int64

		start := r.URL.Query().Get("start")
		if start == "" {
			// Default start time to the start of today (00:00:00)
			startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
			startMillis = toUnixMillis(startTime)
		} else {
			// Parse the incoming start time and convert it to Unix milliseconds
			if parsedTime, err := time.ParseInLocation("2006-01-02T15:04", start, loc); err == nil {
				startMillis = toUnixMillis(parsedTime)
			}
		}

		end := r.URL.Query().Get("end")
		if end == "" {
			// Default end time to the current time
			endMillis = toUnixMillis(now)
		} else {
			// Parse the incoming end time and convert it to Unix milliseconds
			if parsedTime, err := time.ParseInLocation("2006-01-02T15:04", end, loc); err == nil {
				endMillis = toUnixMillis(parsedTime)
			}
		}

		logging.Log.Info().Msgf("Reading data with start: %s, end: %s", start, end)

		commands := collector.GetAllCommandsForPeriod(startMillis, endMillis)
		processes := collector.GetAllProcessesForPeriod(startMillis, endMillis)

		timeProcesses, err := collector.GetTopProcessesAndMetrics(startMillis, endMillis)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Serialize the commands and processes to JSON strings
		commandsJSON, err := json.Marshal(commands)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		processesJSON, err := json.Marshal(processes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		timeProcessesJSON, err := json.Marshal(timeProcesses)
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

		tmpl.Execute(w, map[string]interface{}{
			"CommandsJSON":      string(commandsJSON),
			"ProcessesJSON":     string(processesJSON),
			"TimeProcessesJSON": string(timeProcessesJSON),
			"StartTime":         start,
			"EndTime":           end,
		})
	})

	http.HandleFunc("/command", func(w http.ResponseWriter, r *http.Request) {

		queryParams := r.URL.Query()

		label := queryParams.Get("label")

		if label == "" {
			http.Error(w, "Pull parameter is required", http.StatusBadRequest)
		}

		loc, _ := time.LoadLocation("Local")
		now := time.Now().In(loc)

		// Function to convert time.Time to Unix milliseconds
		toUnixMillis := func(t time.Time) int64 {
			return t.UnixNano() / int64(time.Millisecond)
		}

		var startMillis, endMillis int64

		start := queryParams.Get("start")
		if start == "" {
			// Default start time to the start of today (00:00:00)
			startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
			startMillis = toUnixMillis(startTime)
		} else {
			// Parse the incoming start time and convert it to Unix milliseconds
			if parsedTime, err := time.ParseInLocation("2006-01-02T15:04", start, loc); err == nil {
				startMillis = toUnixMillis(parsedTime)
			}
		}

		end := queryParams.Get("end")
		if end == "" {
			// Default end time to the current time
			endMillis = toUnixMillis(now)
		} else {
			// Parse the incoming end time and convert it to Unix milliseconds
			if parsedTime, err := time.ParseInLocation("2006-01-02T15:04", end, loc); err == nil {
				endMillis = toUnixMillis(parsedTime)
			}
		}

		commands := collector.GetAllCommandsForCategoryForPeriod(
			label, startMillis, endMillis)

		// Serialize the commands to JSON strings
		commandsJSON, err := json.Marshal(commands)
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

		tmpl.Execute(w, map[string]interface{}{
			"CommandsJSON": string(commandsJSON),
			"StartTime":    start,
			"EndTime":      end,
		})

	})
}
