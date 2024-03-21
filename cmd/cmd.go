package cmd

import (
	"fmt"
	"lda/daemon"
	"lda/logging"
	"lda/resources"
	"lda/shell"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// Verbose Define verbose flag variables
	Verbose bool

	ldaCmd = &cobra.Command{
		Use:   "lda",
		Short: "Command line manager for LDA project.",
		Long: `Command line manager for LDA Project.
        Complete documentation is available at https://devzero.io`,
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
		RunE:  serve,
	}
)

func init() {

	includeShowFlagsForLda(ldaCmd)
	includeShowFlagsForServe(serveCmd)

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

func includeShowFlagsForServe(cmd *cobra.Command) {
	cmd.Flags().StringP("port", "p", "8080", "Port to serve the frontend client")
}

func initLogging() {
	logging.Setup(os.Stdout, Verbose)
}

// Execute is the entry point for the command line
func Execute() {
	if err := ldaCmd.Execute(); err != nil {
		logging.Log.Fatal().Err(err).Msg("Failed to execute main lda command")
		os.Exit(1)
	}
}

func lda(cmd *cobra.Command, _ []string) {
	if err := cmd.Help(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to show help")
	}
}

func start(_ *cobra.Command, _ []string) {
	if err := daemon.StartDaemon(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to start daemon")
	}
}

func stop(_ *cobra.Command, _ []string) {
	if err := daemon.StopDaemon(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to stop daemon")
	}
}

func install(_ *cobra.Command, _ []string) {
	if err := daemon.InstallDaemonConfiguration(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to install daemon configuration")
	}

	if err := shell.InstallShellConfiguration(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to install shell configuration")
	}

	if err := shell.InjectShellSource(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to inject shell source")
	}
}

func uninstall(_ *cobra.Command, _ []string) {
	if err := daemon.DestroyDaemonConfiguration(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to uninstall daemon configuration")
	}

	if err := shell.DeleteShellConfiguration(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to delete shell configuration")
	}

	logging.Log.Info().Msg(`
		Daemon service files and shell configuration deleted successfully, 
		~/.lda directory still holds database file, and your rc file stills has source script.
		If you wish to remove those, delete them manually`)
}

func serve(cmd *cobra.Command, _ []string) error {
	portFlag := cmd.Flag("port").Value

	logging.Log.Info().Msgf("Serving local frontend client on http://localhost:%v", portFlag)

	resources.Serve()

	err := http.ListenAndServe(fmt.Sprintf(":%v", portFlag), nil)
	if err != nil {
		return errors.Wrap(err, "pass a port when calling serve; example: `lda serve -p 8987`")
	}
	return nil
}
