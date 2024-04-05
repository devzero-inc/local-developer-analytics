package shell

import (
	"bytes"
	"embed"
	"fmt"
	"lda/collector"
	"lda/util"
	"os"
	"os/user"
	"path"
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
	// List of supported shells
	supportedShells = []string{"/bin/bash", "/bin/zsh", "/bin/fish"}

	templateSources = map[Type]string{
		Zsh:  "scripts/zsh.sh",
		Bash: "scripts/bash.sh",
		Fish: "scripts/fish.sh",
	}

	sourceScripts = map[Type]string{
		Zsh: `
# LDA shell source
if [ -f "$HOME/.lda/lda.sh" ]; then
    source "$HOME/.lda/lda.sh"
fi`,
		Bash: `
# LDA shell source
if [ -f "$HOME/.lda/lda.sh" ]; then
    source "$HOME/.lda/lda.sh"
fi`,
		Fish: `
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

// Type is the type of the shell that is supported
type Type int

const (
	Bash Type = 0
	Zsh  Type = 1
	Fish Type = 2
	Sh   Type = 3
)

// Shell is the shell configuration
type Shell struct {
	ShellType     Type
	ShellLocation string
	isRoot        bool
	SudoExecUser  *user.User
	logger        zerolog.Logger
	ldaDir        string
	homeDir       string
}

// NewShell creates a new shell configuration
func NewShell(logger zerolog.Logger, isRoot bool, ldaDir string, homeDir string, sudoExecUser *user.User) (*Shell, error) {

	shellType, shellLocation, err := setupShell()

	if err != nil {
		return nil, err
	}

	return &Shell{
		ShellType:     shellType,
		ShellLocation: shellLocation,
		logger:        logger,
		isRoot:        isRoot,
		ldaDir:        ldaDir,
		homeDir:       homeDir,
		SudoExecUser:  sudoExecUser,
	}, nil
}

// setupShell sets the current active shell and location
func setupShell() (Type, string, error) {

	shellLocation := os.Getenv("SHELL")

	return configureShell(shellLocation)
}

func configureShell(shellLocation string) (Type, string, error) {
	shellType := path.Base(shellLocation)

	var shell Type
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

// InstallShellConfiguration installs the shell configuration
func (s *Shell) InstallShellConfiguration() error {

	filePath := filepath.Join(s.ldaDir, ldaScript)

	collectorFilePath := filepath.Join(s.ldaDir, CollectorName)

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

	if err := util.WriteFileAndChown(collectorFilePath, cmdContent.Bytes(), execPermissions, s.SudoExecUser); err != nil {
		s.logger.Err(err).Msg("Failed to write collector files")
		return err
	}

	shellTmplLocation, ok := templateSources[s.ShellType]
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

	if err := util.WriteFileAndChown(filePath, shellContent.Bytes(), execPermissions, s.SudoExecUser); err != nil {
		s.logger.Err(err).Msg("Failed to write shell files")
		return err
	}

	s.logger.Info().Msg("Shell configured successfully")

	return nil
}

// DeleteShellConfiguration removes the shell configuration
func (s *Shell) DeleteShellConfiguration() error {

	filePath := filepath.Join(s.ldaDir, "lda.sh")

	if err := os.Remove(filePath); err != nil {
		s.logger.Err(err).Msg("Failed to remove shell configuration")
		return err
	}

	filePath = filepath.Join(s.ldaDir, "collector.sh")
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
	switch s.ShellType {
	case Zsh:
		shellConfigFile = filepath.Join(s.homeDir, ".zshrc")
	case Bash:
		shellConfigFile = filepath.Join(s.homeDir, ".bashrc")
	case Fish:
		shellConfigFile = filepath.Join(s.homeDir, ".config/fish/config.fish")
	default:
		s.logger.Error().Msg("Unsupported shell")
		return fmt.Errorf("unsupported shell")
	}

	if s.isRoot {
		config, err := promptForShellPath(shellConfigFile)
		if err != nil {
			return err
		}
		shellConfigFile = config
	}

	source, ok := sourceScripts[s.ShellType]
	if !ok {
		s.logger.Error().Msg("Unsupported shell")
		return fmt.Errorf("unsupported shell")
	}

	s.logger.Debug().Msgf("Shell config file: %s", shellConfigFile)
	// Check if the script is already present to avoid duplicates
	if !util.IsScriptPresent(shellConfigFile, source) {
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
