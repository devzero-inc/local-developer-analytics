package resources

import (
	"embed"
	"encoding/json"
	"lda/collector"
	"lda/logging"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

// Embedding directory
//go:embed views/*
var templateFS embed.FS

type CommandLabelId struct {
	Label string `json:"label"`
	Id    int    `json:"id"`
}

type ChartConfig struct {
	Type    string       `json:"type"`
	Data    ChartData    `json:"data"`
	Options ChartOptions `json:"options"`
}

type ChartData struct {
	Ids      []int           `json:"ids"`
	Labels   []string        `json:"labels"`
	Datasets []ChartDataSets `json:"datasets"`
}

type ChartDataSets struct {
	Label string  `json:"label"`
	Data  []int64 `json:"data"`
}

type ChartOptions struct {
	MaintainAspectRatio bool `json:"maintainAspectRatio"`
	AspectRatio         int  `json:"aspectRatio"`
}

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

		var labels []string
		var ids []int
		var dataPoints []int64

		for _, cmd := range commands {
			labels = append(labels, cmd.Command)
			ids = append(ids, cmd.Id)
			dataPoints = append(dataPoints, cmd.ExecutionTime)
		}

		// Construct the chart data
		chartDataSets := ChartDataSets{
			Label: "Execution Time",
			Data:  dataPoints,
		}

		chartData := ChartData{
			Ids:      ids,
			Labels:   labels,
			Datasets: []ChartDataSets{chartDataSets},
		}

		chartConfig := ChartConfig{
			Type: "pie",
			Data: chartData,
			Options: ChartOptions{
				MaintainAspectRatio: true,
				AspectRatio:         4,
			},
		}

		// Serialize the commands to JSON strings
		chartJSON, err := json.Marshal(chartConfig)
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
			"ChartJSON": string(chartJSON),
			"StartTime": start,
			"EndTime":   end,
		})

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

		command := collector.GetCommandById(i)

		processes := collector.GetAllProcessesForPeriod(command.StartTime, command.EndTime)

		timeProcesses, err := collector.GetTopProcessesAndMetrics(command.StartTime, command.EndTime)
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

		tmpl, err := template.ParseFS(templateFS, "views/overview.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, map[string]interface{}{
			"ProcessesJSON":     string(processesJSON),
			"TimeProcessesJSON": string(timeProcessesJSON),
		})

	})
}
