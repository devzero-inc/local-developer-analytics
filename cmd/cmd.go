package cmd

import (
	"lda/config"
	"lda/daemon"
	"lda/logging"
	"lda/ui"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

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
	logging.Log.Info().Msg("Starting daemon service")

	var cmd *exec.Cmd

	if config.OS == config.Linux {
		cmd = exec.Command(
			"systemctl",
			"--user",
			"start",
			daemon.DaemonServicedName)
	} else if config.OS == config.MacOS {
		path := filepath.Join(
			config.HomeDir,
			daemon.DaemonPlistFilePath,
			daemon.DaemonPlistName)
		cmd = exec.Command(
			"launchctl",
			"load",
			"-w",
			path)
	}

	if err := cmd.Run(); err != nil {
		logging.Log.Err(err).Msg("Failed to start daemon service")
		return
	}

	logging.Log.Info().Msg("Daemon service started successfully")
}

func stop(_ *cobra.Command, _ []string) {
	logging.Log.Info().Msg("Stoping daemon service")

	var cmd *exec.Cmd

	if config.OS == config.Linux {
		cmd = exec.Command(
			"systemctl",
			"--user",
			"stop",
			daemon.DaemonServicedName)
	} else if config.OS == config.MacOS {
		path := filepath.Join(
			config.HomeDir,
			daemon.DaemonPlistFilePath,
			daemon.DaemonPlistName)
		cmd = exec.Command(
			"launchctl",
			"unload",
			"-w",
			path)
	}

	if err := cmd.Run(); err != nil {
		logging.Log.Err(err).Msg("Failed to stop daemon service")
		return
	}

	logging.Log.Info().Msg("Daemon service stoped successfully")
}

func install(_ *cobra.Command, _ []string) {

	logging.Log.Info().Msg("Installing daemon service")

	var filePath string
	var fileContent []byte

	if config.OS == config.Linux {
		filePath = filepath.Join(
			config.HomeDir,
			daemon.DaemonServicedFilePath,
			daemon.DaemonServicedName)
		fileContent = []byte(daemon.DaemonServiced)
	} else if config.OS == config.MacOS {
		filePath = filepath.Join(
			config.HomeDir,
			daemon.DaemonPlistFilePath,
			daemon.DaemonPlistName)
		fileContent = []byte(daemon.DaemonPlist)
	}

	if filePath == "" {
		logging.Log.Error().Msg("Unsupported operating system")
		return
	}

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		logging.Log.Info().Msg("Daemon service file already exists")
		return
	} else if !os.IsNotExist(err) {
		// An error other than "not exists", e.g., permission issues
		logging.Log.Err(err).Msg("Failed to check daemon service file")
		return
	}

	// File does not exist, proceed with writing
	err := os.WriteFile(filePath, fileContent, daemon.DaemonPermission)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to write daemon service file")
		return
	}

	logging.Log.Info().Msg("Daemon service installed successfully")
}

func uninstall(_ *cobra.Command, _ []string) {
	logging.Log.Info().Msg("Uninstalling daemon service")

	var filePath string

	if config.OS == config.Linux {
		filePath = filepath.Join(
			config.HomeDir,
			daemon.DaemonServicedFilePath,
			daemon.DaemonServicedName)
	} else if config.OS == config.MacOS {
		filePath = filepath.Join(
			config.HomeDir,
			daemon.DaemonPlistFilePath,
			daemon.DaemonPlistName)
	}

	if filePath == "" {
		logging.Log.Error().Msg("Unsupported operating system")
		return
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logging.Log.Info().Msg("Daemon service file does not exist")
		return
	} else if err != nil {
		// An error other than "not exists", e.g., permission issues
		logging.Log.Err(err).Msg("Failed to check daemon service file")
		return
	}

	// File exists, proceed with removal
	err := os.Remove(filePath)
	if err != nil {
		logging.Log.Err(err).Msg("Failed to remove daemon service file")
		return
	}

	logging.Log.Info().Msg("Daemon service file removed successfully")
}

func serve(_ *cobra.Command, _ []string) {
	logging.Log.Info().Msg("Serving local frontend client on http://localhost:8080")

	ui.Serve()

	http.ListenAndServe(":8080", nil)
}
