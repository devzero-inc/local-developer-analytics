package config

import (
	"fmt"
	"lda/util"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"

	"github.com/manifoldco/promptui"
)

// ShellType is the type of the shell that is supported
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

// GetUserConfig returns the user configuration
func GetUserConfig() (*user.User, bool) {
	isRoot := os.Geteuid() == 0

	sudoUser := os.Getenv("SUDO_USER")

	if isRoot && sudoUser != "" {
		originalUser, err := user.Lookup(sudoUser)
		if err != nil {
			fmt.Fprintf(SysConfig.ErrOut, "Failed to get user that executed sudo: %s\n", err)
			os.Exit(1)
		}

		return originalUser, isRoot
	}

	return nil, isRoot
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

// GetOS returns the operating system
func GetOS() OSType {
	OSName = runtime.GOOS
	var osType OSType
	switch OSName {
	case "linux":
		osType = Linux
	case "darwin":
		osType = MacOS
	default:
		// TODO: check if this will work on WSL, maybe it will?
		fmt.Fprintf(SysConfig.ErrOut, "Unsupported operating system: %s\n", OSName)
		os.Exit(1)
	}

	return osType
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

// GetHomeDir returns the user home directory
func GetHomeDir(isRoot bool, user *user.User) string {
	if isRoot && user != nil {
		return user.HomeDir
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(SysConfig.ErrOut, "Failed to get user home directory: %s\n", err)
		os.Exit(1)
	}

	return home
}

// SetupLdaDir sets the directory for the shell configuration
func SetupLdaDir() {
	dir := filepath.Join(HomeDir, ".lda")
	if err := util.CreateDirAndChown(dir, os.ModePerm, SudoExecUser); err != nil {
		fmt.Fprintf(SysConfig.ErrOut, "Failed to create LDA home directory: %s\n", err)
		os.Exit(1)
	}

	LdaDir = dir
}

// GetLdaDir returns the directory for the shell configuration
func GetLdaDir(homeDir string, user *user.User) string {
	dir := filepath.Join(homeDir, ".lda")
	if err := util.CreateDirAndChown(dir, os.ModePerm, user); err != nil {
		fmt.Fprintf(SysConfig.ErrOut, "Failed to create LDA home directory: %s\n", err)
		os.Exit(1)
	}

	return dir
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

// GetLdaBinaryPath returns the path to the lda binary
func GetLdaBinaryPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return exePath, nil
}

// GetShell sets the current active shell and location
func GetShell() (ShellType, string, error) {

	shellLocation := os.Getenv("SHELL")

	return configureShell(shellLocation)
}

func configureShell(shellLocation string) (ShellType, string, error) {
	shellType := path.Base(shellLocation)

	var shell ShellType
	switch shellType {
	case "bash":
		shell = Bash
	case "zsh":
		shell = Zsh
	case "fish":
		shell = Fish
		// TODO: consider supporting "sh" and "ash" as well.
	default:
		shellLocation, err := promptForShellType()
		if err != nil {
			return -1, "", err
		}
		return configureShell(shellLocation)
	}

	return shell, shellLocation, nil
}

// promptForShellPath prompts the user to confirm the detected shell path or input a new one.
func promptForShellType() (string, error) {

	supportedShells := []string{"/bin/bash", "/bin/zsh", "/bin/fish"}

	prompt := promptui.Select{
		Label: "We detected an unsupported shell, often this could happen because the script was run as sudo. Currently, we support the following shells. Please select one:",
		Items: supportedShells,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}
