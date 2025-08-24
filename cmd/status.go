/*
Copyright © 2025 Takeru Furuse
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/tsuperis3112/pmdr/internal/client"
	"github.com/tsuperis3112/pmdr/internal/display"
)

// StatusCmd represents the status command
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Shows the current status of the timer",
	Long:  `Shows the current status of the timer (e.g., session type, remaining time).`,
	Run: func(cmd *cobra.Command, args []string) {
		reply, err := client.Status()
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		display.Status(reply)
	},
}
