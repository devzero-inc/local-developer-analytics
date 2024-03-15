package main

import (
	"lda/cmd"
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
}

func main() {
	cmd.Execute()
}
