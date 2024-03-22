package cmd

import (
	"lda/client"
	"lda/collector"
	"lda/config"
	"lda/logging"

	"github.com/spf13/cobra"
)

var (
	collectCmd = &cobra.Command{
		Use:   "collect",
		Short: "Collect command and system information",
		Long:  `Collect and process command and system information.`,

		Run: collect,
	}
)

func collect(_ *cobra.Command, _ []string) {

	var grpcClient *client.Client
	var err error
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
		}
	}

	collectorInstance := collector.NewCollector(
		collector.SocketPath,
		grpcClient,
	)
	collectorInstance.Collect()
}
