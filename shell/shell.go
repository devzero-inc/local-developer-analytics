package shell

import (
	"bytes"
	"embed"
	"fmt"
	"lda/collector"
	"lda/config"
	"lda/util"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog"
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

// Config is the configuration for the shell
type Config struct {
	ShellType     config.ShellType
	ShellLocation string
	IsRoot        bool
	SudoExecUser  *user.User
	LdaDir        string
	HomeDir       string
}

// Shell is the shell configuration
type Shell struct {
	logger zerolog.Logger
	Config *Config
}

// NewShell creates a new shell configuration
func NewShell(config *Config, logger zerolog.Logger) (*Shell, error) {

	return &Shell{
		logger: logger,
		Config: config,
	}, nil
}

// InstallShellConfiguration installs the shell configuration
func (s *Shell) InstallShellConfiguration() error {

	filePath := filepath.Join(s.Config.LdaDir, ldaScript)

	collectorFilePath := filepath.Join(s.Config.LdaDir, CollectorName)

	cmdTmpl, err := template.ParseFS(templateFS, CollectorScript)
	if err != nil {
		s.logger.Err(err).Msg("Failed to parse collector template")
		return err
	}

	var cmdContent bytes.Buffer
	if err := cmdTmpl.Execute(&cmdContent, map[string]interface{}{
		"SocketPath": collector.SocketPath,
	}); err != nil {
		s.logger.Err(err).Msg("Failed to execute cmd template")
		return err
	}

	if err := util.WriteFileAndChown(collectorFilePath, cmdContent.Bytes(), execPermissions, s.Config.SudoExecUser); err != nil {
		s.logger.Err(err).Msg("Failed to write collector files")
		return err
	}

	shellTmplLocation, ok := templateSources[s.Config.ShellType]
	if !ok {
		s.logger.Error().Msg("Unsupported shell")
		return fmt.Errorf("unsupported shell located")
	}

	shellTmpl, err := template.ParseFS(templateFS, shellTmplLocation)
	if err != nil {
		s.logger.Err(err).Msg("Failed to parse shell template")
		return err
	}

	var shellContent bytes.Buffer
	if err := shellTmpl.Execute(&shellContent, map[string]interface{}{
		"CommandScriptPath": collectorFilePath,
	}); err != nil {
		s.logger.Err(err).Msg("Failed to execute shell template")
		return err
	}

	if err := util.WriteFileAndChown(filePath, shellContent.Bytes(), execPermissions, s.Config.SudoExecUser); err != nil {
		s.logger.Err(err).Msg("Failed to write shell files")
		return err
	}

	s.logger.Info().Msg("Shell configured successfully")

	return nil
}

// DeleteShellConfiguration removes the shell configuration
func (s *Shell) DeleteShellConfiguration() error {

	filePath := filepath.Join(s.Config.LdaDir, "lda.sh")

	if err := os.Remove(filePath); err != nil {
		s.logger.Err(err).Msg("Failed to remove shell configuration")
		return err
	}

	filePath = filepath.Join(s.Config.LdaDir, "collector.sh")
	if err := os.Remove(filePath); err != nil {
		s.logger.Err(err).Msg("Failed to remove shell configuration")
		return err
	}

	s.logger.Info().Msg("Shell configuration removed successfully")

	return nil
}

// InjectShellSource injects the shell source
func (s *Shell) InjectShellSource() error {
	s.logger.Info().Msg("Installing shell source")

	var shellConfigFile string
	switch s.Config.ShellType {
	case config.Zsh:
		shellConfigFile = filepath.Join(s.Config.HomeDir, ".zshrc")
	case config.Bash:
		shellConfigFile = filepath.Join(s.Config.HomeDir, ".bashrc")
	case config.Fish:
		shellConfigFile = filepath.Join(s.Config.HomeDir, ".config/fish/config.fish")
	default:
		s.logger.Error().Msg("Unsupported shell")
		return fmt.Errorf("unsupported shell")
	}

	if s.Config.IsRoot {
		autoupdate := os.Getenv("LDA_AUTO_UPDATE_CONFIG")
		// if user has set autoupdate in the env var, lets stick to installing in the default path for that shell
		// so this logic will not prompt the user for shell info if autoupdate is set
		// TODO support all shells for auto-update mechanisms
		if !strings.EqualFold(autoupdate, "true") {
			conf, err := promptForShellPath(shellConfigFile)
			if err != nil {
				return err
			}
			shellConfigFile = conf
		}
	}

	source, ok := sourceScripts[s.Config.ShellType]
	if !ok {
		s.logger.Error().Msg("Unsupported shell")
		return fmt.Errorf("unsupported shell")
	}

	s.logger.Debug().Msgf("Shell config file: %s", shellConfigFile)
	// Check if the script is already present to avoid duplicates
	if !util.IsScriptPresent(shellConfigFile, "LDA shell source") {
		if err := util.AppendToFile(shellConfigFile, source); err != nil {
			s.logger.Error().Msg("Failed to append to the file")
			return err
		}
	}

	s.logger.Info().Msg("Shell source injected successfully")

	return nil
}

// promptForShellPath uses prompt to ask the user to confirm or enter a new shell path.
func promptForShellPath(detectedShellPath string) (string, error) {
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("We will try to inject this into your shell located at the path: %s. If this is not your shell path, input the path to the shell where we can inject the source; if it is, just press Enter", detectedShellPath),
		Default:   detectedShellPath,
		AllowEdit: true,
		Validate: func(input string) error {
			// TODO: check if path exists
			return nil
		},
	}

	result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	result = strings.TrimSpace(result)
	return result, nil
}
