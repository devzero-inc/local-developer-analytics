package cmd

import (
	"lda/collector"

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
	collectorInstance := collector.NewCollector(collector.SocketPath)
	collectorInstance.Collect()
}
