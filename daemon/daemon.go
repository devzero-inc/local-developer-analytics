package daemon

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/spf13/afero"
	"html/template"
	"lda/config"
	"lda/logging"
	"os/exec"
	"path/filepath"
)

const (
	PlistFilePath               = "Library/LaunchAgents"
	PlistName                   = "lda.plist"
	UserServicedFilePath        = ".config/systemd/user"
	RootServicedFilePath        = "/etc/systemd/system"
	ServicedName                = "lda.service"
	LinuxDaemonTemplateLocation = "services/lda.service"
	MacOSDaemonTemplateLocation = "services/lda.plist"
	ServicePermission           = 0644
	DirPermission               = 0755
)

// Embedding scripts directory
//
//go:embed services/*
var templateFS embed.FS

// InstallDaemonConfiguration installs the daemon service configuration
func InstallDaemonConfiguration() error {
	logging.Log.Info().Msg("Installing daemon service...")

	filePath, templatePath, err := buildConfigurationPath()
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFS(templateFS, templatePath)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to parse config template")
		return err
	}

	var content bytes.Buffer
	if err := tmpl.Execute(&content, map[string]interface{}{
		"BinaryPath": config.ExePath,
		"Shell":      config.ShellLocation,
		"Home":       config.HomeDir,
	}); err != nil {
		logging.Log.Err(err).Msg("Failed to execute daemon template")
		return err
	}

	if err := config.Fs.MkdirAll(filepath.Dir(filePath), DirPermission); err != nil {
		logging.Log.Err(err).Msg("Failed to create directories for daemon files")
		return fmt.Errorf("failed to create directories for daemon files: %w", err)
	}

	if err := afero.WriteFile(config.Fs, filePath, content.Bytes(), ServicePermission); err != nil {
		logging.Log.Err(err).Msg("Failed to write daemon files")
		return fmt.Errorf("failed to write daemon files: %w", err)
	}

	logging.Log.Info().Msg("Daemon service installed successfully")

	return nil
}

// DestroyDaemonConfiguration removes the daemon service configuration
func DestroyDaemonConfiguration() error {
	logging.Log.Info().Msg("Uninstalling daemon service...")

	filePath, _, err := buildConfigurationPath()
	if err != nil {
		return err
	}

	if err := config.Fs.Remove(filePath); err != nil {
		logging.Log.Err(err).Msg("Failed to remove daemon service file")
		return err
	}

	logging.Log.Info().Msg("Daemon service file removed successfully")

	return nil
}

// StartDaemon starts the daemon service
func StartDaemon() error {
	logging.Log.Info().Msg("Starting daemon service...")

	switch config.OS {
	case config.Linux:
		if err := startLinuxDaemon(); err != nil {
			logging.Log.Err(err).Msg("Failed to start daemon service")
			return err
		}
	case config.MacOS:
		if err := startMacOSDaemon(); err != nil {
			logging.Log.Err(err).Msg("Failed to start daemon service")
			return err
		}
	default:
		logging.Log.Error().Msg("Unsupported operating system")
		return fmt.Errorf("unsupported operating system")
	}

	logging.Log.Info().Msg("Daemon service started successfully")

	return nil
}

// StopDaemon stops the daemon service
func StopDaemon() error {
	logging.Log.Info().Msg("Stopping daemon service")

	switch config.OS {
	case config.Linux:
		if err := stopLinuxDaemon(); err != nil {
			logging.Log.Err(err).Msg("Failed to stop daemon service")
			return err
		}
	case config.MacOS:
		if err := stopMacOSDaemon(); err != nil {
			logging.Log.Err(err).Msg("Failed to stop daemon service")
			return err
		}
	default:
		logging.Log.Error().Msg("Unsupported operating system")
		return fmt.Errorf("unsupported operating system")
	}

	logging.Log.Info().Msg("Daemon service stopped successfully")

	return nil
}

// ReloadDaemon signals the daemon to reload its configuration.
func ReloadDaemon() error {
	logging.Log.Info().Msg("Reloading daemon service...")

	switch config.OS {
	case config.Linux:
		return reloadLinuxDaemon()
	case config.MacOS:
		return reloadMacOSDaemon()
	default:
		logging.Log.Error().Msg("Unsupported operating system for reload")
		return fmt.Errorf("unsupported operating system")
	}
}

// reloadLinuxDaemon reloads the daemon service on Linux using systemctl.
func reloadLinuxDaemon() error {
	cmd := exec.Command("systemctl", "--user", "reload", ServicedName)
	if config.IsRoot {
		cmd = exec.Command("systemctl", "reload", ServicedName)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		logging.Log.Err(err).Msgf("Failed to reload daemon service: %v", stderr.String())
		return fmt.Errorf("failed to reload daemon service: %v", stderr.String())
	}

	logging.Log.Info().Msg("Daemon service reloaded successfully")
	return nil
}

// reloadMacOSDaemon reloads the daemon service on macOS
func reloadMacOSDaemon() error {
	stopErr := stopMacOSDaemon()
	if stopErr != nil {
		return stopErr
	}
	startErr := startMacOSDaemon()
	if startErr != nil {
		return startErr
	}

	logging.Log.Info().Msg("Daemon service reloaded successfully")
	return nil
}

// startLinuxDaemon starts the daemon service on Linux
func startLinuxDaemon() error {
	if !checkLogindService() && !config.IsRoot {
		return fmt.Errorf("logind service is not available, and you need to be root to enable the daemon service, or enable logind service manually")
	}

	enableCmd := exec.Command("systemctl", "--user", "enable", ServicedName)
	if config.IsRoot {
		enableCmd = exec.Command("systemctl", "enable", ServicedName)
	}
	var stderr bytes.Buffer
	enableCmd.Stderr = &stderr

	if err := enableCmd.Run(); err != nil {
		return fmt.Errorf("failed to enable daemon service: %v", stderr.String())
	}

	cmd := exec.Command("systemctl", "--user", "start", ServicedName)
	if config.IsRoot {
		cmd = exec.Command("systemctl", "start", ServicedName)
	}
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start daemon service: %v", stderr.String())
	}

	return nil
}

// startMacOSDaemon starts the daemon service on macOS
func startMacOSDaemon() error {
	path := filepath.Join(config.HomeDir, PlistFilePath, PlistName)
	cmd := exec.Command("launchctl", "load", "-w", path)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start daemon service: %v", stderr.String())
	}

	return nil
}

// stopLinuxDaemon stops the daemon service on Linux
func stopLinuxDaemon() error {
	cmd := exec.Command("systemctl", "--user", "stop", ServicedName)
	if config.IsRoot {
		cmd = exec.Command("systemctl", "stop", ServicedName)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop daemon service: %v", stderr.String())
	}

	return nil
}

// stopMacOSDaemon stops the daemon service on macOS
func stopMacOSDaemon() error {
	path := filepath.Join(config.HomeDir, PlistFilePath, PlistName)
	cmd := exec.Command("launchctl", "unload", "-w", path)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop daemon service: %v", stderr.String())
	}

	return nil
}

// buildConfigurationPath builds the path to the daemon service configuration file, and the template
func buildConfigurationPath() (string, string, error) {
	var filePath string
	var templateLocation string

	switch config.OS {
	case config.Linux:
		servicePath := filepath.Join(config.HomeDir, UserServicedFilePath)

		if !checkLogindService() && !config.IsRoot {
			logging.Log.Info().
				Msg("You need to be root to install the daemon service, or enable logind service manually")
			return "", "", fmt.Errorf("logind service is not available")
		}

		if config.IsRoot {
			servicePath = RootServicedFilePath
		}

		filePath = filepath.Join(servicePath, ServicedName)
		templateLocation = LinuxDaemonTemplateLocation
	case config.MacOS:
		filePath = filepath.Join(config.HomeDir, PlistFilePath, PlistName)
		templateLocation = MacOSDaemonTemplateLocation
	default:
		return "", "", fmt.Errorf("unsupported operating system")
	}

	return filePath, templateLocation, nil
}

// Check if logind service is available on the system, because there are Linux systems where
// users deliberately disable logind service, and in such cases, the daemon service will not work
// on a user level, and we have to force sudo usage
func checkLogindService() bool {
	var stderr bytes.Buffer
	cmd := exec.Command("systemctl", "is-enabled", "systemd-logind.service")
	cmd.Stderr = &stderr
	err := cmd.Run()
	logindStatus := stderr.String()

	if err != nil || logindStatus == "masked" || logindStatus == "disabled" {
		logging.Log.Info().Msgf("Logind service is not available, status: %s", logindStatus)
		return false
	}

	return true
}
