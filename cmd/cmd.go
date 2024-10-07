package cmd

import (
	"fmt"
	"io"
	"lda/config"
	"lda/daemon"
	"lda/database"
	"lda/job"
	"lda/logging"
	"lda/resources"
	"lda/shell"
	"lda/user"
	"lda/util"
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

	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Print current configuration",
		Long:  `Display current configuration for LDA Project.`,
		RunE:  displayConfig,
	}

	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Serve local client",
		Long:  `Serve local frontend client for LDA Project.`,
		RunE:  serve,
	}
)

const (
	// days are amount of days that old data will be retained
	days = 5
	// hours are amount of hours that ticker in job will use to run cleanup job
	hours = 24
)

func init() {

	includeShowFlagsForLda(ldaCmd)
	includeShowFlagsForServe(serveCmd)
	includeShowFlagsForInstall(installCmd)

	cobra.OnInitialize(setupConfig)

	ldaCmd.AddCommand(versionCmd)
	ldaCmd.AddCommand(collectCmd)
	ldaCmd.AddCommand(startCmd)
	ldaCmd.AddCommand(stopCmd)
	ldaCmd.AddCommand(installCmd)
	ldaCmd.AddCommand(uninstallCmd)
	ldaCmd.AddCommand(serveCmd)
	ldaCmd.AddCommand(reloadCmd)
	ldaCmd.AddCommand(configCmd)
}

func includeShowFlagsForLda(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbosity")
}

func includeShowFlagsForServe(cmd *cobra.Command) {
	cmd.Flags().StringP("port", "p", "8080", "Port to serve the frontend client")
}

func includeShowFlagsForInstall(cmd *cobra.Command) {
	cmd.Flags().BoolP("auto-credentials", "a", false, "Try to automatically generate the credentails")
	cmd.Flags().BoolP("workspace", "w", false, "Is collection executed in a DevZero workspace")
}

func setupConfig() {

	// setting up the system configuration
	config.SetupSysConfig()

	// Setup afero FS layer
	util.SetupFS()

	sudoExecUser, isRoot, err := config.GetUserConfig()
	if err != nil {
		fmt.Fprintf(config.SysConfig.ErrOut, "Failed to get user configuration: %s\n", err)
		os.Exit(1)
	}
	osConf, osName, err := config.GetOS()
	if err != nil {
		fmt.Fprintf(config.SysConfig.ErrOut, "Failed to get OS: %s\n", err)
		os.Exit(1)
	}
	homeDir, err := config.GetHomeDir(isRoot, sudoExecUser)
	if err != nil {
		fmt.Fprintf(config.SysConfig.ErrOut, "Failed to get home directory: %s\n", err)
		os.Exit(1)
	}
	ldaDir, err := config.GetLdaDir(homeDir, sudoExecUser)
	if err != nil {
		fmt.Fprintf(config.SysConfig.ErrOut, "Failed to get LDA directory: %s\n", err)
		os.Exit(1)
	}
	exePath, err := config.GetLdaBinaryPath()
	if err != nil {
		fmt.Fprintf(config.SysConfig.ErrOut, "Failed to get executable path: %s\n", err)
		os.Exit(1)
	}

	// setting up optional application configuration
	config.SetupConfig(ldaDir, sudoExecUser)

	// setup database and run migrations
	database.Setup(ldaDir, sudoExecUser)
	database.RunMigrations()

	// setting up the Logger
	// TODO: consider adding verbose levels
	if config.AppConfig.Debug || Verbose {
		logging.Setup(os.Stdout, true)
	} else {
		logging.Setup(io.Discard, false)
	}

	// Configure default user globals
	user.Conf = &user.Config{
		Os:      int64(osConf),
		OsName:  osName,
		HomeDir: homeDir,
		LdaDir:  ldaDir,
		IsRoot:  isRoot,
		ExePath: exePath,
		User:    sudoExecUser,
	}

	// run cleanup job
	job.Cleanup(hours, days)
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

	user.ConfigureUserSystemInfo(user.Conf)

	daemonConf := &daemon.Config{
		ExePath:       user.Conf.ExePath,
		HomeDir:       user.Conf.HomeDir,
		IsRoot:        user.Conf.IsRoot,
		Os:            config.OSType(user.Conf.Os),
		SudoExecUser:  user.Conf.User,
		ShellLocation: user.Conf.ShellLocation,
	}
	dmn := daemon.NewDaemon(daemonConf, logging.Log)

	fmt.Fprintln(config.SysConfig.Out, "Reloading LDA daemon...")
	if err := dmn.ReloadDaemon(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to reload daemon")
		return errors.Wrap(err, "failed to reload LDA daemon")
	}
	fmt.Fprintln(config.SysConfig.Out, "Reloading LDA daemon finished.")
	return nil
}

func start(_ *cobra.Command, _ []string) error {

	user.ConfigureUserSystemInfo(user.Conf)

	daemonConf := &daemon.Config{
		ExePath:       user.Conf.ExePath,
		HomeDir:       user.Conf.HomeDir,
		IsRoot:        user.Conf.IsRoot,
		Os:            config.OSType(user.Conf.Os),
		SudoExecUser:  user.Conf.User,
		ShellLocation: user.Conf.ShellLocation,
	}
	dmn := daemon.NewDaemon(daemonConf, logging.Log)

	fmt.Fprintln(config.SysConfig.Out, "Starting LDA daemon...")
	if err := dmn.StartDaemon(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to start daemon")
		return errors.Wrap(err, "failed to start LDA daemon")
	}
	fmt.Fprintln(config.SysConfig.Out, "LDA daemon started.")
	return nil
}

func stop(_ *cobra.Command, _ []string) error {

	user.ConfigureUserSystemInfo(user.Conf)

	daemonConf := &daemon.Config{
		ExePath:       user.Conf.ExePath,
		HomeDir:       user.Conf.HomeDir,
		IsRoot:        user.Conf.IsRoot,
		Os:            config.OSType(user.Conf.Os),
		SudoExecUser:  user.Conf.User,
		ShellLocation: user.Conf.ShellLocation,
	}
	dmn := daemon.NewDaemon(daemonConf, logging.Log)

	fmt.Fprintln(config.SysConfig.Out, "Stopping LDA daemon...")
	if err := dmn.StopDaemon(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to stop daemon")
		return errors.Wrap(err, "failed to stop LDA daemon")
	}
	fmt.Fprintln(config.SysConfig.Out, "LDA daemon stopped.")
	return nil
}

func install(cmd *cobra.Command, _ []string) error {

	autoCredentials, err := cmd.Flags().GetBool("auto-credentials")
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to get auto-credentials flag")
		return errors.Wrap(err, "failed to get auto-credentials flag")
	}

	isWorkspace, err := cmd.Flags().GetBool("workspace")
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to get workspace flag")
		return errors.Wrap(err, "failed to get workspace flag")
	}
	user.ConfigureUserSystemInfo(user.Conf)

	daemonConf := &daemon.Config{
		ExePath:        user.Conf.ExePath,
		HomeDir:        user.Conf.HomeDir,
		IsRoot:         user.Conf.IsRoot,
		Os:             config.OSType(user.Conf.Os),
		SudoExecUser:   user.Conf.User,
		ShellLocation:  user.Conf.ShellLocation,
		AutoCredential: autoCredentials,
		IsWorkspace:    isWorkspace,
	}
	dmn := daemon.NewDaemon(daemonConf, logging.Log)

	fmt.Fprintln(config.SysConfig.Out, "Installing LDA daemon...")
	if err := dmn.InstallDaemonConfiguration(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to install daemon configuration")
		return errors.Wrap(err, "failed to install LDA daemon configuration file")
	}

	shellConfig := &shell.Config{
		ShellType:     config.ShellType(user.Conf.ShellType),
		ShellLocation: user.Conf.ShellLocation,
		IsRoot:        user.Conf.IsRoot,
		SudoExecUser:  user.Conf.User,
		LdaDir:        user.Conf.LdaDir,
		HomeDir:       user.Conf.HomeDir,
	}

	shl, err := shell.NewShell(shellConfig, logging.Log)

	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to setup shell")
		os.Exit(1)
	}

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

	user.ConfigureUserSystemInfo(user.Conf)

	daemonConf := &daemon.Config{
		ExePath:       user.Conf.ExePath,
		HomeDir:       user.Conf.HomeDir,
		IsRoot:        user.Conf.IsRoot,
		Os:            config.OSType(user.Conf.Os),
		SudoExecUser:  user.Conf.User,
		ShellLocation: user.Conf.ShellLocation,
	}
	dmn := daemon.NewDaemon(daemonConf, logging.Log)

	shellConfig := &shell.Config{
		ShellType:     config.ShellType(user.Conf.ShellType),
		ShellLocation: user.Conf.ShellLocation,
		IsRoot:        user.Conf.IsRoot,
		SudoExecUser:  user.Conf.User,
		LdaDir:        user.Conf.LdaDir,
		HomeDir:       user.Conf.HomeDir,
	}
	shl, err := shell.NewShell(shellConfig, logging.Log)

	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to setup shell")
		os.Exit(1)
	}

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

func displayConfig(_ *cobra.Command, _ []string) error {
	conf, err := user.GetConfig()
	if err != nil {
		logging.Log.Error().Err(err).Msg("Failed to get os config")
		return errors.Wrap(err, "failed to get os config, please run 'lda install' first")
	}

	fmt.Fprintf(config.SysConfig.Out, "Current configuration: %+v \n", conf)

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
