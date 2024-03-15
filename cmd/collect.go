package cmd

import (
	"fmt"

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
	fmt.Println("LDA v0.0.1")
}
