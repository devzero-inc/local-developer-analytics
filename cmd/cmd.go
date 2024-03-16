package cmd

import (
	"lda/daemon"
	"lda/logging"
	"lda/resources"
	"lda/shell"
	"net/http"
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

	installCmd = &cobra.Command{
		Use:   "install",
		Short: "Install daemon runner",
		Long:  `Install daemon runner for LDA Project.`,
		Run:   install,
	}

	uninstallCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "Install daemon runner",
		Long:  `Uninstall daemon runner for LDA Project.`,
		Run:   uninstall,
	}

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start daemon runner",
		Long:  `Start daemon runner for LDA Project.`,
		Run:   start,
	}

	stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop daemon runner",
		Long:  `Stop daemon runner for LDA Project.`,
		Run:   stop,
	}

	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Serve local client",
		Long:  `Serve local frontend client for LDA Project.`,
		Run:   serve,
	}
)

func init() {

	includeShowFlagsForLda(ldaCmd)

	cobra.OnInitialize(initLogging)

	ldaCmd.AddCommand(versionCmd)
	ldaCmd.AddCommand(collectCmd)
	ldaCmd.AddCommand(startCmd)
	ldaCmd.AddCommand(stopCmd)
	ldaCmd.AddCommand(installCmd)
	ldaCmd.AddCommand(uninstallCmd)
	ldaCmd.AddCommand(serveCmd)
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

func start(_ *cobra.Command, _ []string) {
	daemon.StartDaemon()
}

func stop(_ *cobra.Command, _ []string) {
	daemon.StopDaemon()
}

func install(_ *cobra.Command, _ []string) {

	shell.InitShellConfiguration()

	daemon.InitDaemonConfiguration()

	shell.InjectShellSource()
}

func uninstall(_ *cobra.Command, _ []string) {
	daemon.DestroyDaemonConfiguration()
}

func serve(_ *cobra.Command, _ []string) {
	logging.Log.Info().Msg("Serving local frontend client on http://localhost:8080")

	resources.Serve()

	http.ListenAndServe(":8080", nil)
}
