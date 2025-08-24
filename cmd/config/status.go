/*
Copyright © 2025 Takeru Furuse
*/
package config

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// StatusCmd represents the status command
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the path of the current configuration file",
	Long:  `Show the path of the current configuration file. If no config file is used, it will show (no config).`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile := viper.ConfigFileUsed()
		if configFile != "" {
			slog.Info(configFile)
		} else {
			slog.Info("(no config)")
		}
	},
}
