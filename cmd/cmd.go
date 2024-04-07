package cmd

import (
	"fmt"
	"io"
	"lda/config"
	"lda/daemon"
	"lda/database"
	"lda/logging"
	"lda/resources"
	"lda/shell"
	"lda/user"
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
		RunE:  install,
	}

	uninstallCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "Install daemon runner",
		Long:  `Uninstall daemon runner for LDA Project.`,
		RunE:  uninstall,
	}

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start daemon runner",
		Long:  `Start daemon runner for LDA Project.`,
		RunE:  start,
	}

	stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop daemon runner",
		Long:  `Stop daemon runner for LDA Project.`,
		RunE:  stop,
	}

	reloadCmd = &cobra.Command{
		Use:   "reload",
		Short: "Reload daemon runner",
		Long:  `Reload daemon runner for LDA Project.`,
		RunE:  reload,
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

	cobra.OnInitialize(setupConfig)

	ldaCmd.AddCommand(versionCmd)
	ldaCmd.AddCommand(collectCmd)
	ldaCmd.AddCommand(startCmd)
	ldaCmd.AddCommand(stopCmd)
	ldaCmd.AddCommand(installCmd)
	ldaCmd.AddCommand(uninstallCmd)
	ldaCmd.AddCommand(serveCmd)
	ldaCmd.AddCommand(reloadCmd)
}

func includeShowFlagsForLda(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbosity")
}

func includeShowFlagsForServe(cmd *cobra.Command) {
	cmd.Flags().StringP("port", "p", "8080", "Port to serve the frontend client")
}

func setupConfig() {

	// setting up the system configuration
	config.SetupSysConfig()

	// setting up the user permission level
	config.SetupUserConfig()

	// setting up the operating system
	config.SetupOS()

	// setting up the LDA binary path
	config.SetupLdaBinaryPath()

	// setting up the home directory
	config.SetupHomeDir()

	// setting up the LDA directory
	config.SetupLdaDir()

	// setting up optional application configuration
	config.SetupConfig()

	// setting up the Logger
	// TODO: consider adding verbose levels
	if config.AppConfig.Debug || Verbose {
		logging.Setup(os.Stdout, true)
	} else {
		logging.Setup(io.Discard, false)
	}

	// setup database and run migrations
	database.Setup()
	database.RunMigrations()

	user.CheckAndConfigureGlobals()
}

func setupShell() *shell.Shell {

	shellType, shellLocation, err := config.GetShell()
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to setup shell")
		os.Exit(1)
	}

	shellConfig := &shell.Config{
		ShellType:     shellType,
		ShellLocation: shellLocation,
		IsRoot:        config.IsRoot,
		SudoExecUser:  config.SudoExecUser,
		LdaDir:        config.LdaDir,
		HomeDir:       config.HomeDir,
	}

	shl, err := shell.NewShell(shellConfig, logging.Log)

	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to setup shell")
		os.Exit(1)
	}

	return shl
}

// setupDaemon initializes the daemon runner with basic configuration
func setupDaemon() *daemon.Daemon {

	daemonConf := &daemon.Config{
		ExePath:      config.ExePath,
		HomeDir:      config.HomeDir,
		IsRoot:       config.IsRoot,
		Os:           config.OS,
		SudoExecUser: config.SudoExecUser,
	}
	return daemon.NewDaemon(daemonConf, logging.Log)
}

// Execute is the entry point for the command line
func Execute() {
	if err := ldaCmd.Execute(); err != nil {
		logging.Log.Err(err).Msg("Failed to execute main lda command")
		os.Exit(1)
	}
}

func lda(cmd *cobra.Command, _ []string) {
	if err := cmd.Help(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to show help")
	}
}

func reload(_ *cobra.Command, _ []string) error {

	dmn := setupDaemon()

	fmt.Fprintln(config.SysConfig.Out, "Reloading LDA daemon...")
	if err := dmn.ReloadDaemon(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to reload daemon")
		return errors.Wrap(err, "failed to reload LDA daemon")
	}
	fmt.Fprintln(config.SysConfig.Out, "Reloading LDA daemon finished.")
	return nil
}

func start(_ *cobra.Command, _ []string) error {

	dmn := setupDaemon()

	fmt.Fprintln(config.SysConfig.Out, "Starting LDA daemon...")
	if err := dmn.StartDaemon(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to start daemon")
		return errors.Wrap(err, "failed to start LDA daemon")
	}
	fmt.Fprintln(config.SysConfig.Out, "LDA daemon started.")
	return nil
}

func stop(_ *cobra.Command, _ []string) error {

	dmn := setupDaemon()

	fmt.Fprintln(config.SysConfig.Out, "Stopping LDA daemon...")
	if err := dmn.StopDaemon(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to stop daemon")
		return errors.Wrap(err, "failed to stop LDA daemon")
	}
	fmt.Fprintln(config.SysConfig.Out, "LDA daemon stopped.")
	return nil
}

func install(_ *cobra.Command, _ []string) error {

	dmn := setupDaemon()

	fmt.Fprintln(config.SysConfig.Out, "Installing LDA daemon...")
	if err := dmn.InstallDaemonConfiguration(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to install daemon configuration")
		return errors.Wrap(err, "failed to install LDA daemon configuration file")
	}

	shl := setupShell()

	if err := shl.InstallShellConfiguration(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to install shell configuration")
		return errors.Wrap(err, "failed to install LDA shell configuration files")
	}

	if err := shl.InjectShellSource(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to inject shell source")
		return errors.Wrap(err, "failed to inject LDA shell source")
	}

	fmt.Fprintln(config.SysConfig.Out, "LDA daemon installed successfully.")
	return nil
}

func uninstall(_ *cobra.Command, _ []string) error {

	shl := setupShell()
	dmn := setupDaemon()

	fmt.Fprintln(config.SysConfig.Out, "Uninstalling LDA daemon...")
	if err := dmn.DestroyDaemonConfiguration(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to uninstall daemon configuration")
		return errors.Wrap(err, "failed to uninstall LDA daemon configuration file")
	}

	if err := shl.DeleteShellConfiguration(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to delete shell configuration")
		return errors.Wrap(err, "failed to delete LDA shell configuration files")
	}

	fmt.Fprintln(config.SysConfig.Out, `Daemon service files and shell configuration deleted successfully, 
		~/.lda directory still holds database file, and your rc file stills has source script.
		If you wish to remove those, delete them manually`)

	return nil
}

func serve(cmd *cobra.Command, _ []string) error {
	portFlag := cmd.Flag("port").Value

	fmt.Fprintf(config.SysConfig.Out, "Serving local frontend client on http://localhost:%v\n", portFlag)

	resources.Serve()

	err := http.ListenAndServe(fmt.Sprintf(":%v", portFlag), nil)
	if err != nil {
		return errors.Wrap(err, "pass a port when calling serve; example: `lda serve -p 8987`")
	}
	return nil
}
