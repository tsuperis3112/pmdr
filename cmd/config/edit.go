/*
Copyright © 2025 Takeru Furuse
*/
package config

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tsuperis3112/pmdr/internal/config"
)

// EditCmd represents the edit command
var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the configuration file",
	Long:  `Open the current configuration file in the default editor.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configFile := viper.ConfigFileUsed()
		if configFile == "" {
			// If no config file is in use, default to the home config path
			configDir, configPath, err := config.GetDefaultConfigPaths()
			if err != nil {
				return fmt.Errorf("failed to get default config path: %w", err)
			}
			configFile = configPath

			// Ensure the directory exists
			if _, err := os.Stat(configDir); os.IsNotExist(err) {
				if err := os.MkdirAll(configDir, 0755); err != nil {
					return fmt.Errorf("failed to create config directory: %w", err)
				}
			}
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim" // default editor
		}

		editorCmd := exec.Command(editor, configFile)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		if err := editorCmd.Run(); err != nil {
			return fmt.Errorf("failed to open editor: %w", err)
		}

		return nil
	},
}
