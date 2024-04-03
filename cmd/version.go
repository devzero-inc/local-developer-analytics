package cmd

import (
	"fmt"
	"lda/config"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version number of LDA.",
		Long:  `Display LDA version number.`,

		Run: version,
	}
)

func version(_ *cobra.Command, _ []string) {
	fmt.Fprintf(config.SysConfig.Out, "LDA v%s\n", config.Version)
}
