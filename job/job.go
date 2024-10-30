package job

import (
	"time"

	"github.com/devzero-inc/local-developer-analytics/collector"
	"github.com/devzero-inc/local-developer-analytics/process"
)

// Cleanup job that will run in background and every 'hours' try to run the ticker
// and delete process and commands older than 'days'
func Cleanup(hours int, days int) {
	// ticker to run cleanup every n hours
	ticker := time.NewTicker(time.Duration(hours) * time.Hour)

	go func() {
		for {
			select {
			case <-ticker.C:
				collector.DeleteCommandsByDays(days)
				process.DeleteProcessesByDays(days)
			}
		}
	}()
}
