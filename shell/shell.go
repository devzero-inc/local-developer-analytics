package shell

import (
	"bufio"
	"bytes"
	"embed"
	"lda/collector"
	"lda/config"
	"lda/logging"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	shellDir        = ".lda"
	execPermissions = 0755
)

// Embedding scripts directory
//go:embed scripts/*
var templateFS embed.FS

func InitShellConfiguration() {
	logging.Log.Info().Msg("Installing shell configuration")

	dir := filepath.Join(config.HomeDir, shellDir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil && !os.IsExist(err) {
		logging.Log.Err(err).Msg("Failed to create shell configuration directory")
	}

	var filePath string
	var shellScriptLocation string

	switch config.Shell {
	case config.Zsh:
		filePath = filepath.Join(dir, "lda.sh")
		shellScriptLocation = "scripts/zsh.sh"
	default:
		logging.Log.Error().Msg("Unsupported shell")
		return
	}

	collectorFilePath := filepath.Join(dir, "collector.sh")

	cmdTmpl, err := template.ParseFS(templateFS, "scripts/collector.sh")
	if err != nil {
		logging.Log.Err(err).Msg("Failed to parse collector template")
		return
	}

	var cmdContent bytes.Buffer
	if err := cmdTmpl.Execute(&cmdContent, map[string]interface{}{
		"SocketPath": collector.SocketPath,
	}); err != nil {
		logging.Log.Err(err).Msg("Failed to execute cmd template")
		return
	}

	if err := os.WriteFile(collectorFilePath, cmdContent.Bytes(), execPermissions); err != nil {
		logging.Log.Err(err).Msg("Failed to write collector files")
		return
	}

	shellTempl, err := template.ParseFS(templateFS, shellScriptLocation)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to parse shell template")
		return
	}

	var shellContent bytes.Buffer
	if err := shellTempl.Execute(&shellContent, map[string]interface{}{
		"CommandScriptPath": collectorFilePath,
	}); err != nil {
		logging.Log.Err(err).Msg("Failed to execute shell template")
		return
	}

	if err := os.WriteFile(filePath, shellContent.Bytes(), execPermissions); err != nil {
		logging.Log.Err(err).Msg("Failed to write shell files")
		return
	}

	logging.Log.Info().Msg("Shell configured successfully")
}

func InjectShellSource() {

	logging.Log.Info().Msg("Installing shell source")

	var shellConfigFile string

	switch config.Shell {
	case config.Zsh:
		shellConfigFile = filepath.Join(config.HomeDir, ".zshrc")

	default:
		logging.Log.Error().Msg("Unsupported shell")
		return
	}

	script := `
# LDA shell source
if [ -f "$HOME/.lda/lda.sh" ]; then
    source "$HOME/.lda/lda.sh"
fi`

	// Check if the script is already present to avoid duplicates
	if !isScriptPresent(shellConfigFile, script) {
		if err := appendToFile(shellConfigFile, script); err != nil {
			return
		}
	}

	logging.Log.Info().Msg("Shell source injected successfully")
}

func isScriptPresent(filePath, script string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), ".lda/lda.sh") {
			return true
		}
	}
	return false
}

func appendToFile(filePath, content string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return err
	}
	return nil
}
