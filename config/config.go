package config

import (
	"github.com/spf13/viper"
	"lda/logging"
)

// Config config definition
type Config struct {
	Debug                     bool   `mapstructure:"debug"`
	ProcessInterval           int    `mapstructure:"process_interval"`
	CommandInterval           int    `mapstructure:"command_interval"`
	CommandIntervalMultiplier int    `mapstructure:"command_interval_multiplier"`
	RemoteCollection          bool   `mapstructure:"remote_collection"`
	ServerHost                string `mapstructure:"server_host"`
	ServerPort                int    `mapstructure:"server_port"`
	ExcludeRegex              string `mapstructure:"exclude_regex"`
	SecureConnection          bool   `mapstructure:"secure_connection"`
	CertFile                  string `mapstructure:"cert_file"`
}

// AppConfig is the global configuration instance
var AppConfig *Config

// SetupConfig initialize the configuration instance
func SetupConfig() {

	logging.Log.Debug().Msg("Setting up application configuration")

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(LdaDir)

	// Set default configuration values
	var config = &Config{
		Debug:                     false,
		RemoteCollection:          false,
		ProcessInterval:           120,
		CommandInterval:           1,
		CommandIntervalMultiplier: 5,
	}

	if err := viper.ReadInConfig(); err != nil {
		logging.Log.Debug().Err(err).Msg("Failed to read config file")
	}

	if err := viper.Unmarshal(config); err != nil {
		logging.Log.Debug().Err(err).Msg("Failed to unmarshal config file")
	}

	logging.Log.Debug().Msgf("Config loaded: %+v", config)

	AppConfig = config
}
