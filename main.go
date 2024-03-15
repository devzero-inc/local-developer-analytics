package main

import (
	"lda/cmd"
	"lda/config"
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
		true,
	)

	// setting up the operating system
	config.SetupOs()

	// setting up the home directory
	config.SetupHomeDir()
}

func main() {
	cmd.Execute()
}
