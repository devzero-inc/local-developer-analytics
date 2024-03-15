package collector

import (
	"lda/database"
	"lda/logging"
)

type Command struct {
	Id            int    `json:"id" db:"id"`
	PID           int    `json:"pid" db:"pid"`
	Command       string `json:"command" db:"command"`
	User          string `json:"user" db:"user"`
	Directory     string `json:"directory" db:"directory"`
	ExecutionTime string `json:"executionTime" db:"execution_time"`
	StartTime     string `json:"startTime" db:"start_time"`
	EndTime       string `json:"endTime" db:"end_time"`
}

func GetAllCommands() {
	var commands []Command
	if err := database.DB.Select(&commands, "SELECT * FROM commands"); err != nil {
		logging.Log.Err(err).Msg("Failed to get all commands")
	}
}

func InsertCommand(command Command) {
	query := `INSERT INTO commands (command, user, directory, execution_time, start_time, end_time)
VALUES (:command, :user, :directory, :execution_time, :start_time, :end_time)`

	_, err := database.DB.NamedExec(query, command)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to insert command")
	}
}
