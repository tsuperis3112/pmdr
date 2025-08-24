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

// resumeCmd represents the resume command
var ResumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resumes a paused session",
	Long:  `Resumes a paused session. The timer will continue from where it left off.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.Resume(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		slog.Info("Pomodoro session resumed.")
	},
}
