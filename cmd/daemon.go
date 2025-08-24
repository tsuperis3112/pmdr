/*
Copyright © 2025 Takeru Furuse
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/tsuperis3112/pmdr/internal/daemon"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:    "daemon",
	Short:  "Run the pmdr daemon process",
	Long:   `This command starts the pmdr daemon process that runs in the background.`,
	Hidden: true, // This makes the command hidden from the help message
	Run: func(cmd *cobra.Command, args []string) {
		if err := daemon.Run(); err != nil {
			slog.Error("Daemon failed", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(daemonCmd)
}
