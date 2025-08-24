package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tsuperis3112/pmdr/cmd/config"
	configInternal "github.com/tsuperis3112/pmdr/internal/config"
	"github.com/tsuperis3112/pmdr/internal/logging"
)

var (
	cfgFile  string
	logLevel string
	logPath  string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "pmdr",
	Short: "A simple Pomodoro timer for your terminal",
	Long: `pmdr is a minimalist, yet powerful Pomodoro Technique timer
designed to run as a background daemon, helping you stay focused
and productive right from your command line.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		level := slog.LevelInfo // default
		if err := level.UnmarshalText([]byte(logLevel)); err != nil {
			level = slog.LevelInfo
		}
		logging.Init(level, logPath)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Initialize sub-packages
	config.Initialize()

	// Add subcommands
	RootCmd.AddCommand(StartCmd)
	RootCmd.AddCommand(StatusCmd)
	RootCmd.AddCommand(PauseCmd)
	RootCmd.AddCommand(ResumeCmd)
	RootCmd.AddCommand(StopCmd)
	RootCmd.AddCommand(config.Cmd)

	// Persistent flags
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/pmdr/config.yaml)")
	RootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level (debug, info, warn, error)")
	RootCmd.PersistentFlags().StringVar(&logPath, "log-path", "", "log file path (default is stderr)")

	// Viper binding
	vip := viper.GetViper()
	_ = vip.BindPFlag("log.level", RootCmd.PersistentFlags().Lookup("log-level"))
	_ = vip.BindPFlag("log.path", RootCmd.PersistentFlags().Lookup("log-path"))
}

func initConfig() {
	foundCfgFile, err := configInternal.FindConfigFile(cfgFile)
	cobra.CheckErr(err)

	if foundCfgFile != "" {
		viper.SetConfigFile(foundCfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()

	vip := viper.GetViper()

	// Set log level and path from viper if not set by flags
	if logLevel == "info" && vip.IsSet("log.level") {
		logLevel = vip.GetString("log.level")
	}
	if logPath == "" && vip.IsSet("log.path") {
		logPath = vip.GetString("log.path")
	}
}
