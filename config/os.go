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

// GetUserConfig returns the user configuration
func GetUserConfig() (*user.User, bool, error) {
	isRoot := os.Geteuid() == 0

	sudoUser := os.Getenv("SUDO_USER")

	if isRoot && sudoUser != "" {
		originalUser, err := user.Lookup(sudoUser)
		if err != nil {
			return nil, isRoot, err
		}

		return originalUser, isRoot, nil
	}

	return nil, isRoot, nil
}

// GetOS returns the operating system
func GetOS() (OSType, string, error) {
	osName := runtime.GOOS
	var osType OSType
	switch osName {
	case "linux":
		osType = Linux
	case "darwin":
		osType = MacOS
	default:
		// TODO: check if this will work on WSL, maybe it will?
		return -1, "", fmt.Errorf("unsupported operating system")
	}

	return osType, osName, nil
}

// GetHomeDir returns the user home directory
func GetHomeDir(isRoot bool, user *user.User) (string, error) {
	if isRoot && user != nil {
		return user.HomeDir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return home, nil
}

// GetLdaDir returns the directory for the shell configuration
func GetLdaDir(homeDir string, user *user.User) (string, error) {
	dir := filepath.Join(homeDir, ".lda")
	if err := util.CreateDirAndChown(dir, os.ModePerm, user); err != nil {
		return "", err
	}

	return dir, nil
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

	ShellLocation = os.Getenv("SHELL")
	logging.Log.Info().Msgf("Trying to determine the shell: %s", ShellLocation)

	shellType := path.Base(ShellLocation)

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
