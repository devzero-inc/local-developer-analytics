package collector

import (
	"lda/config"
	"lda/database"
	"lda/logging"
	"regexp"
)

type Command struct {
	Id            int    `json:"id" db:"id"`
	PID           int    `json:"pid" db:"pid"`
	Category      string `json:"category" db:"category"`
	Command       string `json:"command" db:"command"`
	User          string `json:"user" db:"user"`
	Directory     string `json:"directory" db:"directory"`
	ExecutionTime int64  `json:"executionTime" db:"execution_time"`
	StartTime     int64  `json:"startTime" db:"start_time"`
	EndTime       int64  `json:"endTime" db:"end_time"`
}

func GetAllCommands() []Command {
	var commands []Command
	if err := database.DB.Select(&commands, "SELECT * FROM commands"); err != nil {
		logging.Log.Err(err).Msg("Failed to get all commands")
	}

	return commands
}

func GetCommandById(id int64) Command {
	var command Command
	query := `SELECT * FROM commands WHERE id = ?`

	err := database.DB.Get(&command, query, id)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to get command by id")
	}

	return command
}

func GetAllCommandsForPeriod(start int64, end int64) []Command {
	var commands []Command

	query := `SELECT id, category, SUM(execution_time) AS execution_time 
              FROM commands 
              WHERE start_time >= ? AND start_time <= ? 
              GROUP BY category 
              ORDER BY category ASC, SUM(execution_time) DESC`

	err := database.DB.Select(&commands, query, start, end)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to get aggregated commands with start and end times")
	}

	return commands
}

func GetAllCommandsForCategoryForPeriod(category string, start int64, end int64) []Command {
	var commands []Command

	query := `SELECT id, category, command, SUM(execution_time) AS execution_time 
              FROM commands 
              WHERE category = ? AND start_time >= ? AND start_time <= ? 
              GROUP BY command 
              ORDER BY command ASC, SUM(execution_time) DESC`

	err := database.DB.Select(&commands, query, category, start, end)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to get aggregated commands with start and end times")
	}

	return commands
}

//func GetAllCommandsForPeriod(start, end string) []Command {
//	var commands []Command
//
//	query := `SELECT * FROM commands WHERE start_time >= ? AND end_time <= ? ORDER BY start_time ASC`
//
//	err := database.DB.Select(&commands, query, start, end)
//	if err != nil {
//		logging.Log.Err(err).Msg("Failed to get commands with start and end times")
//	}
//
//	return commands
//}

func InsertCommand(command Command) {
	query := `INSERT INTO commands (category, command, user, directory, execution_time, start_time, end_time)
	VALUES (:category, :command, :user, :directory, :execution_time, :start_time, :end_time)`

	_, err := database.DB.NamedExec(query, command)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to insert command")
	}
}

func ParseCommand(command string) string {

	// TODO: there might be some other cases as well like: watch, time etc
	// we might need to figure out how to handle them
	var pattern = regexp.MustCompile(`^(?:sudo|nohup)?\s*(?:\./|/usr/bin/|/bin/|/usr/local/bin/)?([^/ ]+?)(?:\s|$)`)

	matches := pattern.FindStringSubmatch(command)
	if len(matches) > 1 {
		return matches[1]
	}

	return command
}

func IsCommandAcceptable(command string) bool {

	if config.AppConfig.ExcludeRegex != "" {
		logging.Log.Debug().Msgf("Checking if command is acceptable: %s", command)
		var pattern = regexp.MustCompile(config.AppConfig.ExcludeRegex)
		return !pattern.MatchString(command)
	}

	return true
}
