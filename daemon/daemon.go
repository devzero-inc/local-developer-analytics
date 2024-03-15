package daemon

import (
	"bytes"
	"embed"
	"html/template"
	"lda/config"
	"lda/logging"
	"os"
	"path/filepath"
)

const (
	DaemonPlistFilePath    = "Library/LaunchAgents"
	DaemonPlistName        = "lda.plist"
	DaemonServicedFilePath = ".config/systemd/user"
	DaemonServicedName     = "lda.service"
	DaemonPermission       = 0644
)

// Embedding scripts directory
//go:embed configs/*
var templateFS embed.FS

func InitDaemonConfiguration() {

	logging.Log.Info().Msg("Installing daemon service")

	var filePath string
	var configLocation string

	if config.OS == config.Linux {
		filePath = filepath.Join(
			config.HomeDir,
			DaemonServicedFilePath,
			DaemonServicedName)
		configLocation = "configs/lda.service"
	} else if config.OS == config.MacOS {
		filePath = filepath.Join(
			config.HomeDir,
			DaemonPlistFilePath,
			DaemonPlistName)
		configLocation = "configs/lda.plist"
	}

	if filePath == "" {
		logging.Log.Error().Msg("Unsupported operating system")
		return
	}

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		logging.Log.Info().Msg("Daemon service file already exists")
		return
	} else if !os.IsNotExist(err) {
		// An error other than "not exists", e.g., permission issues
		logging.Log.Err(err).Msg("Failed to check daemon service file")
		return
	}

	shellTempl, err := template.ParseFS(templateFS, configLocation)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to parse config template")
		return
	}

	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	var content bytes.Buffer
	if err := shellTempl.Execute(&content, map[string]interface{}{
		"BinaryPath": exePath,
	}); err != nil {
		logging.Log.Err(err).Msg("Failed to execute daemon template")
		return
	}

	if err := os.WriteFile(filePath, content.Bytes(), DaemonPermission); err != nil {
		logging.Log.Err(err).Msg("Failed to write daemon files")
		return
	}

	logging.Log.Info().Msg("Daemon service installed successfully")
}
