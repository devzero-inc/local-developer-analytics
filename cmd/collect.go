package cmd

import (
	"lda/collector"
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
	logging.Log.Info().Msg("Collecting command and system information")

	collector.Collect()

	logging.Log.Info().Msg("Collection stoped")
}
