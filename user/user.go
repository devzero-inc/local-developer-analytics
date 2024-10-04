package user

import (
	"database/sql"
	"errors"
	"fmt"
	"lda/collector"
	"lda/config"
	"lda/database"
	"lda/logging"
	"lda/util"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

// TODO: set up localization engine for these strings.
const (
	YesUpdate string = "Yes, update to new configuration"
	NoKeep    string = "No, keep the existing configuration"
)

var Conf *Config

// Config is the basic configuration for the system
type Config struct {
	Id int64 `json:"id" db:"id"`
	// OS is the operating system
	Os int64 `json:"os" db:"os"`
	// OsName is the operating system name
	OsName string `json:"os_name" db:"os_name"`
	// HomeDir is the user home directory
	HomeDir string `json:"home_dir" db:"home_dir"`
	// LdaDir is the home LDA directory where all configurations are stored.
	LdaDir string `json:"lda_dir" db:"lda_dir"`
	// IsRoot is a value to check if the user is root
	IsRoot bool `json:"is_root" db:"is_root"`
	// ExePath is the path to the lda binary
	ExePath string `json:"exe_path" db:"exe_path"`
	// ShellType is the type of the shell
	ShellType int64 `json:"shell_type" db:"shell_type"`
	// ShellLocation is the location of the shell
	ShellLocation string `json:"shell_location" db:"shell_location"`
	// User is the user that executed the command (if sudo)
	User *user.User `json:"-" db:"-"`
}

// GetConfig fetches Config used to configure the system
func GetConfig() (*Config, error) {
	var osConfig Config
	query := `SELECT * FROM config LIMIT 1`

	if err := database.DB.Get(&osConfig, query); err != nil {
		logging.Log.Err(err).Msg("Failed to get os config")
		return nil, err
	}

	return &osConfig, nil
}

// InsertConfig inserts Config used to configure the system
func InsertConfig(osConfig Config) error {
	query := `INSERT INTO config (os, os_name, home_dir, lda_dir, is_root, exe_path, shell_type, shell_location) 
			  VALUES (:os, :os_name, :home_dir, :lda_dir, :is_root, :exe_path, :shell_type, :shell_location)`

	_, err := database.DB.NamedExec(query, osConfig)

	return err
}

// UpdateConfig updates an existing Config record in the database
func UpdateConfig(osConfig Config) error {
	query := `UPDATE config SET 
                os = :os, 
                os_name = :os_name, 
                home_dir = :home_dir, 
                lda_dir = :lda_dir, 
                is_root = :is_root, 
                exe_path = :exe_path,
                shell_type = :shell_type,
                shell_location = :shell_location
              WHERE id = :id`

	_, err := database.DB.NamedExec(query, osConfig)

	return err
}

// ConfigureUserSystemInfo configures the user system information and prompts the user to update the configuration if necessary.
func ConfigureUserSystemInfo(currentConf *Config) {

	// Retrieve the existing configuration from the database.
	existingConf, err := GetConfig()
	if err != nil && err != sql.ErrNoRows {
		logging.Log.Err(err).Msg("Failed to get os config")
		fmt.Fprintf(config.SysConfig.ErrOut, "Failed to get os config: %s\n", err)
		os.Exit(1)
	}

	// If config exists, compare it with the current configuration.
	if existingConf != nil {
		logging.Log.Debug().Msg("Existing config found, comparing...")

		hasDiff, diffs := CompareConfig(existingConf, currentConf)
		if hasDiff {
			logging.Log.Info().Msg("Configuration drift detected")
			diffMsg := strings.Join(diffs, "\n")
			fmt.Fprintf(config.SysConfig.Out, "Differences detected in configuration:\n%s\n", diffMsg)

			// allow for non-user interrupted flow
			var result string
			autoupdate := os.Getenv("LDA_AUTO_UPDATE_CONFIG")
			if strings.EqualFold(autoupdate, "true") {
				result = YesUpdate
			}

			// if no env vars, prompt the user
			if result == "" {
				// Prompt user to choose between old and new config
				prompt := promptui.Select{
					Label: "Configuration drift detected. Do you want to update the configuration to the new settings?",
					Items: []string{YesUpdate, NoKeep},
				}

				_, result, err = prompt.Run()
				if err != nil {
					logging.Log.Error().Err(err).Msg("Prompt failed")
					fmt.Fprintf(config.SysConfig.ErrOut, "Prompt failed: %s\n", err)
					os.Exit(1)
				}
			}

			if result == YesUpdate {
				shellType, shellLocation, err := config.GetShell()
				if err != nil {
					logging.Log.Error().Err(err).Msg("Failed to setup shell")
					os.Exit(1)
				}

				currentConf.ShellType = int64(shellType)
				currentConf.ShellLocation = shellLocation

				currentConf.Id = existingConf.Id
				if err := UpdateConfig(*currentConf); err != nil {
					logging.Log.Error().Err(err).Msg("Failed to update configuration")
					fmt.Fprintf(config.SysConfig.ErrOut, "Failed to update configuration: %s\n", err)
					os.Exit(1)
				}
				logging.Log.Info().Msg("Configuration updated to current settings.")
				Conf = currentConf
			} else {
				existingConf.User = currentConf.User
				Conf = existingConf
			}
		} else {
			logging.Log.Debug().Msg("No configuration drift detected.")
			existingConf.User = currentConf.User
			Conf = existingConf
		}

		return
	}

	logging.Log.Debug().Msg("No config found, creating new one")

	shellType, shellLocation, err := config.GetShell()
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to setup shell")
		os.Exit(1)
	}

	currentConf.ShellType = int64(shellType)
	currentConf.ShellLocation = shellLocation

	logging.Log.Debug().Msgf("Shell config: %+v", currentConf)

	if err := InsertConfig(*currentConf); err != nil {
		fmt.Fprintf(config.SysConfig.ErrOut, "Failed to insert os config: %s\n", err)
		os.Exit(1)
	}

	logging.Log.Debug().Msg("Config inserted")

	Conf = currentConf

	logging.Log.Debug().Msgf("Config: %+v", config.SysConfig)
}

// CompareConfig compares two configuration structs and returns a boolean indicating if there are differences,
// along with a slice of strings describing what those differences are.
func CompareConfig(existingConf, currentConf *Config) (bool, []string) {
	var diffs []string

	if existingConf.Os != currentConf.Os {
		diffs = append(diffs, fmt.Sprintf("Operating System changed from %d to %d", existingConf.Os, currentConf.Os))
	}
	if existingConf.OsName != currentConf.OsName {
		diffs = append(diffs, fmt.Sprintf("Operating System Name changed from %s to %s", existingConf.OsName, currentConf.OsName))
	}
	if existingConf.HomeDir != currentConf.HomeDir {
		diffs = append(diffs, fmt.Sprintf("Home Directory changed from %s to %s", existingConf.HomeDir, currentConf.HomeDir))
	}
	if existingConf.LdaDir != currentConf.LdaDir {
		diffs = append(diffs, fmt.Sprintf("LDA Directory changed from %s to %s", existingConf.LdaDir, currentConf.LdaDir))
	}
	if existingConf.IsRoot != currentConf.IsRoot {
		diffs = append(diffs, fmt.Sprintf("IsRoot status changed from %t to %t", existingConf.IsRoot, currentConf.IsRoot))
	}
	if existingConf.ExePath != currentConf.ExePath {
		diffs = append(diffs, fmt.Sprintf("Executable Path changed from %s to %s", existingConf.ExePath, currentConf.ExePath))
	}

	return len(diffs) > 0, diffs
}

// GetStoragePath returns the path to the devzero storage directory based on the operating system
func GetStoragePath(os config.OSType, home string) (string, error) {
	switch os {
	case config.MacOS:
		return filepath.Join(home, "Library", "Application Support", "devzero"), nil
	case config.Linux:
		return filepath.Join(home, ".local", "share", "devzero"), nil
	default:
		return "", errors.New("unsupported os")
	}
}

// ReadDZWorkspaceConfig reads the DevZero workspace configuration
func ReadDZWorkspaceConfig() (collector.AuthConfig, error) {
	const (
		devzeroConfigPath    = "/etc/devzero/configs"
		devzeroTeamFile      = "DEVZERO_TEAM_ID"
		devzeroWorkspaceFile = "DEVZERO_WORKSPACE_ID"
	)

	userId := ""
	teamId := ""
	workspaceId := ""

	teamPath := filepath.Join(devzeroConfigPath, devzeroTeamFile)
	if util.FileExists(teamPath) {
		data, err := os.ReadFile(teamPath)
		if err == nil && len(data) > 0 {
			teamId = string(data)
		}
	}

	workspacePath := filepath.Join(devzeroConfigPath, devzeroWorkspaceFile)
	if util.FileExists(workspacePath) {
		data, err := os.ReadFile(workspacePath)
		if err == nil && len(data) > 0 {
			workspaceId = string(data)
		}
	}

	return collector.AuthConfig{
		UserID:      userId,
		TeamID:      teamId,
		WorkspaceID: workspaceId,
	}, nil
}

// ReadDZCliConfig reads the DevZero workspace configuration
func ReadDZCliConfig(path string) (collector.AuthConfig, error) {
	const (
		localUserFile = "user_id.txt"
		localTeamFile = "team_id.txt"
	)

	userId := ""
	teamId := ""

	localUserPath := filepath.Join(path, localUserFile)
	if util.FileExists(localUserPath) {
		data, err := os.ReadFile(localUserPath)
		if err == nil && len(data) > 0 {
			userId = string(data)
		}
	}

	localTeamPath := filepath.Join(path, localTeamFile)
	if util.FileExists(localTeamPath) {
		data, err := os.ReadFile(localTeamPath)
		if err == nil && len(data) > 0 {
			teamId = string(data)
		}
	}

	return collector.AuthConfig{
		UserID: userId,
		TeamID: teamId,
	}, nil
}
