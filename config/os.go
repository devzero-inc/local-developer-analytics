package config

import (
	"lda/logging"
	"os"
	"runtime"
)

type OSType int

const (
	Linux OSType = 0
	MacOS OSType = 1
)

var (
	// OS is the operating system
	OS OSType
	// HomeDir is the user home directory
	HomeDir string
)

// SetupOs determine the operating system
func SetupOs() {
	switch runtime.GOOS {
	case "linux":
		logging.Log.Info().Msg("Running on Linux")
		OS = Linux
	case "darwin":
		logging.Log.Info().Msg("Running on macOS")
		OS = MacOS
	default:
		logging.Log.Fatal().Msg("Unsupported operating system")
		os.Exit(1)
	}
}

// SetHomeDir sets the user home directory
func SetupHomeDir() {
	home, err := os.UserHomeDir()
	if err != nil {
		logging.Log.Fatal().Err(err).Msg("Failed to get user home directory")
		os.Exit(1)
	}
	HomeDir = home
}
