package config

import (
	"lda/logging"
	"os"
	"path/filepath"
	"runtime"
)

type ShellType int

const (
	Bash ShellType = 0
	Zsh  ShellType = 1
	Fish ShellType = 2
	Sh   ShellType = 3
)

type OSType int

const (
	Linux OSType = 0
	MacOS OSType = 1
)

var (
	// OS is the operating system
	OS OSType
	// Shell is the active shell
	Shell ShellType
	// HomeDir is the user home directory
	HomeDir string
	// LdaDir is the directory.
	LdaDir string
	// IsRoot is a value to check if the user is root
	IsRoot bool
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

// SetupShell sets the current active shell
func SetupShell() {

	shell := os.Getenv("SHELL")
	logging.Log.Info().Msgf("Trying to determine the shell: %s", shell)

	switch shell {
	case "/bin/bash":
		logging.Log.Info().Msg("Running bash shell")
		Shell = Bash
	case "/bin/zsh":
		logging.Log.Info().Msg("Running zsh shell")
		Shell = Zsh
	case "/bin/fish":
		logging.Log.Info().Msg("Running fish shell")
		Shell = Fish
	case "/bin/sh":
		logging.Log.Info().Msg("Running sh shell")
		Shell = Sh
	default:
		logging.Log.Fatal().Msg("Unsupported shell")
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

// SetupLdaDir sets the directory for the shell configuration
func SetupLdaDir() {

	dir := filepath.Join(HomeDir, ".lda")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil && !os.IsExist(err) {
		logging.Log.Err(err).Msg("Failed to create shell configuration directory")
	}

	LdaDir = dir
}

// SetupUserConfig sets the user permission level (root or not)
func SetupUserConfig() {
	IsRoot = os.Geteuid() == 0
}
