package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

// ShellType is the type of the shell
type ShellType int

const (
	Bash ShellType = 0
	Zsh  ShellType = 1
	Fish ShellType = 2
	Sh   ShellType = 3
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
	// Shell is the active shell
	Shell ShellType
	// ShellLocation is the shell configuration directory
	ShellLocation string
	// HomeDir is the user home directory
	HomeDir string
	// LdaDir is the home lad directory where all configurations are stored.
	LdaDir string
	// IsRoot is a value to check if the user is root
	IsRoot bool
	// ExePath is the path to the lda binary
	ExePath string
)

// SetupOs determine the operating system
func SetupOs() {
	OSName = runtime.GOOS
	switch OSName {
	case "linux":
		OS = Linux
	case "darwin":
		OS = MacOS
	default:
		// TODO: check if this will work on WSL, maybe it will?
		fmt.Fprint(SysConfig.ErrOut, "Unsupported operating system")
		os.Exit(1)
	}
}

// SetupShell sets the current active shell
func SetupShell() {

	ShellLocation = os.Getenv("SHELL")

	shellType := path.Base(ShellLocation)

	switch shellType {
	case "bash":
		Shell = Bash
	case "zsh":
		Shell = Zsh
	case "fish":
		Shell = Fish
	case "sh":
		Shell = Sh
		// TODO: consider supporting "ash" as well.
	default:
		fmt.Fprint(SysConfig.ErrOut, "Unsupported shell")
		os.Exit(1)
	}

}

// SetupHomeDir sets the user home directory
func SetupHomeDir() {
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

// SetupUserConfig sets the user permission level (root or not)
func SetupUserConfig() {
	IsRoot = os.Geteuid() == 0
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
