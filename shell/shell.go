package shell

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/spf13/afero"
	"lda/collector"
	"lda/config"
	"lda/logging"
	"lda/util"
	"path/filepath"
	"text/template"
)

const (
	ldaScript       = "lda.sh"
	execPermissions = 0755
	CollectorName   = "collector.sh"
	CollectorScript = "scripts/collector.sh"
)

var (
	templateSources = map[config.ShellType]string{
		config.Zsh:  "scripts/zsh.sh",
		config.Bash: "scripts/bash.sh",
		config.Fish: "scripts/fish.sh",
	}

	sourceScripts = map[config.ShellType]string{
		config.Zsh: `
# LDA shell source
if [ -f "$HOME/.lda/lda.sh" ]; then
    source "$HOME/.lda/lda.sh"
fi`,
		config.Bash: `
# LDA shell source
if [ -f "$HOME/.lda/lda.sh" ]; then
    source "$HOME/.lda/lda.sh"
fi`,
		config.Fish: `
# LDA shell source
if test -f "$HOME/.lda/lda.sh"
    source "$HOME/.lda/lda.sh"
end`,
	}

	// Embedding scripts directory
	//
	//go:embed scripts/*
	templateFS embed.FS
)

// InstallShellConfiguration installs the shell configuration
func InstallShellConfiguration() error {

	filePath := filepath.Join(config.LdaDir, ldaScript)

	collectorFilePath := filepath.Join(config.LdaDir, CollectorName)

	cmdTmpl, err := template.ParseFS(templateFS, CollectorScript)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to parse collector template")
		return err
	}

	var cmdContent bytes.Buffer
	if err := cmdTmpl.Execute(&cmdContent, map[string]interface{}{
		"SocketPath": collector.SocketPath,
	}); err != nil {
		logging.Log.Err(err).Msg("Failed to execute cmd template")
		return err
	}

	if err := afero.WriteFile(config.Fs, collectorFilePath, cmdContent.Bytes(), execPermissions); err != nil {
		logging.Log.Err(err).Msg("Failed to write collector files")
		return err
	}

	shellTmplLocation, ok := templateSources[config.Shell]
	if !ok {
		logging.Log.Error().Msg("Unsupported shell")
		return fmt.Errorf("unsupported operating system")
	}

	shellTmpl, err := template.ParseFS(templateFS, shellTmplLocation)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to parse shell template")
		return err
	}

	var shellContent bytes.Buffer
	if err := shellTmpl.Execute(&shellContent, map[string]interface{}{
		"CommandScriptPath": collectorFilePath,
	}); err != nil {
		logging.Log.Err(err).Msg("Failed to execute shell template")
		return err
	}

	if err := afero.WriteFile(config.Fs, filePath, shellContent.Bytes(), execPermissions); err != nil {
		logging.Log.Err(err).Msg("Failed to write shell files")
		return err
	}

	logging.Log.Info().Msg("Shell configured successfully")

	return nil
}

// DeleteShellConfiguration removes the shell configuration
func DeleteShellConfiguration() error {

	filePath := filepath.Join(config.LdaDir, "lda.sh")

	if err := config.Fs.Remove(filePath); err != nil {
		logging.Log.Err(err).Msg("Failed to remove shell configuration")
		return err
	}

	filePath = filepath.Join(config.LdaDir, "collector.sh")
	if err := config.Fs.Remove(filePath); err != nil {
		logging.Log.Err(err).Msg("Failed to remove shell configuration")
		return err
	}

	logging.Log.Info().Msg("Shell configuration removed successfully")

	return nil
}

// InjectShellSource injects the shell source
func InjectShellSource() error {
	logging.Log.Info().Msg("Installing shell source")

	var shellConfigFile string
	switch config.Shell {
	case config.Zsh:
		shellConfigFile = filepath.Join(config.HomeDir, ".zshrc")
	case config.Bash:
		shellConfigFile = filepath.Join(config.HomeDir, ".bashrc")
	case config.Fish:
		shellConfigFile = filepath.Join(config.HomeDir, ".config/fish/config.fish")
	default:
		logging.Log.Error().Msg("Unsupported shell")
		return fmt.Errorf("unsupported shell")
	}

	source, ok := sourceScripts[config.Shell]
	if !ok {
		logging.Log.Error().Msg("Unsupported shell")
		return fmt.Errorf("unsupported shell")
	}

	logging.Log.Debug().Msgf("Shell config file: %s", shellConfigFile)
	// Check if the script is already present to avoid duplicates
	if !util.IsScriptPresent(shellConfigFile, source) {
		if err := util.AppendToFile(shellConfigFile, source); err != nil {
			logging.Log.Error().Msg("Failed to append to the file")
			return err
		}
	}

	logging.Log.Info().Msg("Shell source injected successfully")

	return nil
}
