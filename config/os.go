package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// OSType is the type of the operating system
type OSType int

const (
	Linux OSType = 0
	MacOS OSType = 1
)

var (
	// OS is the operating system
	OS OSType
	// OSName is the name of the operating system
	OSName string
	// HomeDir is the user home directory
	HomeDir string
	// LdaDir is the home LDA directory where all configurations are stored.
	LdaDir string
	// IsRoot is a value to check if the user is root
	IsRoot bool
	// ExePath is the path to the lda binary
	ExePath string
	// SudoExecUser is the user that executed the command (if sudo)
	SudoExecUser *user.User
)

// SetupUserConfig sets the user permission level (root or not)
func SetupUserConfig() {
	IsRoot = os.Geteuid() == 0

	sudoUser := os.Getenv("SUDO_USER")

	if IsRoot && sudoUser != "" {
		originalUser, err := user.Lookup(sudoUser)
		if err != nil {
			fmt.Fprintf(SysConfig.ErrOut, "Failed to get user that executed sudo: %s\n", err)
			os.Exit(1)
		}

		SudoExecUser = originalUser
	}
}

// SetupOS determine the operating system
func SetupOS() {
	OSName = runtime.GOOS
	switch OSName {
	case "linux":
		OS = Linux
	case "darwin":
		OS = MacOS
	default:
		// TODO: check if this will work on WSL, maybe it will?
		fmt.Fprintf(SysConfig.ErrOut, "Unsupported operating system: %s\n", OSName)
		os.Exit(1)
	}
}

// SetupHomeDir sets the user home directory
func SetupHomeDir() {
	if IsRoot && SudoExecUser != nil {
		HomeDir = SudoExecUser.HomeDir
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(SysConfig.ErrOut, "Failed to get user home directory: %s\n", err)
		os.Exit(1)
	}

	HomeDir = home
}

// SetupLdaDir sets the directory for the shell configuration
func SetupLdaDir() {

	dir := filepath.Join(HomeDir, ".lda")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil && !os.IsExist(err) {
		fmt.Fprintf(SysConfig.ErrOut, "Failed to create LDA home directory: %s\n", err)
		os.Exit(1)
	}

	LdaDir = dir
}

// SetupLdaBinaryPath sets the path to the lda binary
func SetupLdaBinaryPath() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(SysConfig.ErrOut, "Failed to find executable path: %s\n", err)
		os.Exit(1)
	}
	ExePath = exePath
}
