package ui

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

		// Default start time to the start of today (00:00:00), formatted for datetime-local input
		start := r.URL.Query().Get("start")
		if start == "" {
			start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).Format("2006-01-02T15:04")
		} else {
			// Parse the incoming start time and reformat it to ensure consistency
			if parsedTime, err := time.Parse(time.RFC3339, start); err == nil {
				start = parsedTime.Format("2006-01-02T15:04")
			}
		}

		// Default end time to the current time, formatted for datetime-local input
		end := r.URL.Query().Get("end")
		if end == "" {
			end = now.Format("2006-01-02T15:04")
		} else {
			// Parse the incoming end time and reformat it to ensure consistency
			if parsedTime, err := time.Parse(time.RFC3339, end); err == nil {
				end = parsedTime.Format("2006-01-02T15:04")
			}
		}

		commands := collector.GetAllCommandsForPeriod(start, end)
		processes := collector.GetAllProcessesForPeriod(start, end)

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

		tmpl, err := template.ParseFS(templateFS, "views/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logging.Log.Info().Msgf("Rendering template with start: %s, end: %s", start, end)

		tmpl.Execute(w, map[string]interface{}{
			"CommandsJSON":  string(commandsJSON),
			"ProcessesJSON": string(processesJSON),
			"StartTime":     start,
			"EndTime":       end,
		})
	})

}
