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

// pauseCmd represents the pause command
var PauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pauses the current session",
	Long:  `Pauses the current session. The timer will stop until resume is called.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.Pause(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		slog.Info("Pomodoro session paused.")
	},
}
