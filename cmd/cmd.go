package cmd

import (
	"lda/logging"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Verbose Define verbose flag variables
	Verbose bool

	ldaCmd = &cobra.Command{
		Use:   "lda",
		Short: "Command line manager for LDA project.",
		Long: `Command line manager for LDA Project.
        Complete documentation is available at http://devzero.io`,

		Run: lda,
	}
)

func init() {

	includeShowFlagsForLda(ldaCmd)

	cobra.OnInitialize(initLogging)

	ldaCmd.AddCommand(versionCmd)
}

func includeShowFlagsForLda(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbosity")
}

func initLogging() {
	logging.Setup(os.Stdout, Verbose)
}

func Execute() {
	if err := ldaCmd.Execute(); err != nil {
		logging.Log.Fatal().Err(err).Msg("Failed to execute main lda command")
		os.Exit(1)
	}
}

func lda(cmd *cobra.Command, _ []string) {
	cmd.Help()
}
