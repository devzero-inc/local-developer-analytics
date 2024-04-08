package cmd

import (
	"lda/client"
	"lda/collector"
	"lda/config"
	"lda/logging"
	"lda/process"
	"lda/user"
	"time"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

var (
	collectCmd = &cobra.Command{
		Use:   "collect",
		Short: "Collect command and system information",
		Long:  `Collect and process command and system information.`,

		RunE: collect,
	}
)

func collect(_ *cobra.Command, _ []string) error {

	var err error

	user.Conf, err = user.GetConfig()
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to get os config")
		return errors.Wrap(err, "failed to get os config, please run 'lda install' first")
	}

	var grpcClient *client.Client
	if config.AppConfig.RemoteCollection {
		logging.Log.Info().Msg("Remote collection is enabled")
		grpcConfig := client.Config{
			Address:          config.AppConfig.ServerAddress,
			SecureConnection: config.AppConfig.SecureConnection,
			CertFile:         config.AppConfig.CertFile,
			Timeout:          60,
		}
		grpcClient, err = client.NewClient(grpcConfig)
		if err != nil {
			logging.Log.Error().Err(err).Msg("Failed to create client")
			return errors.Wrap(err, "failed to create client")
		}
	}

	intervalConfig := collector.IntervalConfig{
		ProcessInterval:           time.Duration(config.AppConfig.ProcessInterval),
		CommandInterval:           time.Duration(config.AppConfig.CommandInterval),
		CommandIntervalMultiplier: time.Duration(config.AppConfig.CommandIntervalMultiplier),
		MaxConcurrentCommands:     config.AppConfig.MaxConcurrentCommands,
	}

	procCol, err := process.NewFactory(logging.Log).Create(config.AppConfig.ProcessCollectionType)
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to create process collector")
		return errors.Wrap(err, "failed to create process collector")
	}

	collectorInstance := collector.NewCollector(
		collector.SocketPath,
		grpcClient,
		logging.Log,
		intervalConfig,
		config.AppConfig.ExcludeRegex,
		procCol,
	)

	collectorInstance.Collect()

	return nil
}
