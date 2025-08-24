/*
Copyright © 2025 Takeru Furuse
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/tsuperis3112/pmdr/internal/client"
)

// stopCmd represents the stop command
var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the timer completely",
	Long:  `Stops the timer completely and terminates the daemon process.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.Stop(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		slog.Info("Pomodoro session stopped.")
	},
}
