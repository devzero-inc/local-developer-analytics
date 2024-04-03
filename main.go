package main

import (
	"lda/cmd"
	"lda/config"
	"lda/database"
	"lda/logging"
	"os"
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

	// setting up the user permission level
	config.SetupUserConfig()

	// setting up optional application configuration
	config.SetupConfig()

	// setting up the LDA binary path
	config.SetupLdaBinaryPath()

	// setup database and run migrations
	database.Setup()
	database.RunMigrations()
}

func main() {
	cmd.Execute()
}
