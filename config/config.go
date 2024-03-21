package config

import (
	"lda/logging"

	"github.com/spf13/viper"
)

// Config config definition
type Config struct {
	Debug            bool   `mapstructure:"debug"`
	RemoteCollection bool   `mapstructure:"remote_collection"`
	ServerHost       string `mapstructure:"server_host"`
	ServerPort       int    `mapstructure:"server_port"`
	ExcludeRegex     string `mapstructure:"exclude_regex"`
}

var AppConfig *Config

// Setup initialize the configuration instance
func SetupConfig() {

	logging.Log.Debug().Msg("Setting up application configuration")

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(LdaDir)

	var config = &Config{
		Debug:            false,
		RemoteCollection: false,
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
