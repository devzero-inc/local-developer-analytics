package main

import (
	"lda/cmd"
	"lda/config"
	"lda/database"
	"lda/logging"
	"os"
)

var (
	Version = "0.0.0"
	Commit  = "xxx"
	Branch  = "undefined"
)

// init function will be called before main
func init() {

	// setting up the Logger
	logging.Setup(
		os.Stdout,
		false,
	)

	// setting up the operating system
	config.SetupOs()

	// setting up the shell
	config.SetupShell()

	// setting up the home directory
	config.SetupHomeDir()

	// setting up the LDA directory
	config.SetupLdaDir()

	// setup database and run migrations
	database.Setup()
	database.RunMigrations()
}

func main() {
	cmd.Execute()
}
