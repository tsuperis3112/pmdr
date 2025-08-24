/*
Copyright © 2025 Takeru Furuse
*/
package config

import (
	"github.com/spf13/cobra"
)

// Cmd represents the config command
var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `Manage pmdr configuration file.`,
}

// Initialize sets up the config command and its subcommands.
func Initialize() {
	Cmd.AddCommand(InitCmd)
	Cmd.AddCommand(EditCmd)
	Cmd.AddCommand(StatusCmd)
}
