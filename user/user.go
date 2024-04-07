package user

import (
	"database/sql"
	"fmt"
	"lda/config"
	"lda/database"
	"lda/logging"
	"os"
)

var userConf *Config

// Config is the basic configuration for the system
type Config struct {
	Id            int64  `json:"id" db:"id"`
	Os            int64  `json:"os" db:"os"`
	OsName        string `json:"os_name" db:"os_name"`
	HomeDir       string `json:"home_dir" db:"home_dir"`
	LdaDir        string `json:"lda_dir" db:"lda_dir"`
	IsRoot        bool   `json:"is_root" db:"is_root"`
	ExePath       string `json:"exe_path" db:"exe_path"`
	ShellType     int64  `json:"shell_type" db:"shell_type"`
	ShellLocation string `json:"shell_location" db:"shell_location"`
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
func InsertConfig(osConfig *Config) error {
	query := `INSERT INTO config (os, os_name, home_dir, lda_dir, is_root, exe_path) 
			  VALUES (:os, :os_name, :home_dir, :lda_dir, :is_root, :exe_path)`

	_, err := database.DB.NamedExec(query, osConfig)

	return err
}

func CheckAndConfigureGlobals() {

	conf, err := GetConfig()
	if err != nil && err != sql.ErrNoRows {
		logging.Log.Err(err).Msg("Failed to get os config")
		fmt.Fprintf(config.SysConfig.ErrOut, "Failed to get os config: %s\n", err)
		os.Exit(1)
	}

	if conf != nil {
		logging.Log.Debug().Msg("Config found")
		userConf = conf
		fmt.Fprintf(config.SysConfig.Out, "Config: %+v\n", conf)
		return
	}

	logging.Log.Debug().Msg("No config found, creating new one")

	shellType, shellLocation, err := config.GetShell()
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to setup shell")
		os.Exit(1)
	}

	osConf := config.GetOS()
	homeDir := config.GetHomeDir(config.IsRoot, config.SudoExecUser)
	ldaDir := config.GetLdaDir(homeDir, config.SudoExecUser)
	exePath, err := config.GetLdaBinaryPath()
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to setup lda binary path")
		os.Exit(1)
	}

	conf = &Config{
		ShellType:     int64(shellType),
		ShellLocation: shellLocation,
		Os:            int64(osConf),
		HomeDir:       homeDir,
		LdaDir:        ldaDir,
		IsRoot:        config.IsRoot,
		ExePath:       exePath,
	}

	fmt.Fprintf(config.SysConfig.Out, "Config: %+v\n", conf)

	if err := InsertConfig(conf); err != nil {
		fmt.Fprintf(config.SysConfig.ErrOut, "Failed to insert os config: %s\n", err)
		os.Exit(1)
	}

	logging.Log.Debug().Msg("Config inserted")

	userConf = conf

	logging.Log.Debug().Msgf("Config: %+v", config.SysConfig)
}
