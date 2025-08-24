/*
Copyright © 2025 Takeru Furuse
*/
package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tsuperis3112/pmdr/internal/config"
)

const defaultConfigTemplate = `# pmdr configuration file
# For more information, see: https://github.com/tsuperis3112/pmdr

# Timer durations (any valid Go time duration string, e.g., "25m", "1h30m")
work_duration: 25m
short_break_duration: 5m
long_break_duration: 15m

# Number of work cycles before a long break
pomo_cycles: 4

# Hooks: execute shell commands on events
hooks:
  # Triggered when a work session finishes
  work:
    # Example for macOS native notification
    # - "osascript -e 'display notification "Work session complete! Time for a break." with title "Pmdr"'"
    # Example for Linux native notification (with libnotify)
    # - "notify-send "Pmdr" "Work session complete! Time for a break."
  # Triggered when a short break session finishes
  short_break:
    # - "osascript -e 'display notification "Break is over! Time for work." with title "Pmdr"'"
  # Triggered when a long break session finishes
  long_break:
    # - "osascript -e 'display notification "Long break is over! Time for work." with title "Pmdr"'"
`

// InitCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a default configuration file",
	Long:  `Create a default configuration file with comments at the default location.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir, configPathYaml, err := config.GetDefaultConfigPaths()
		if err != nil {
			return fmt.Errorf("failed to get default config path: %w", err)
		}
		configPathYml := filepath.Join(configDir, "config.yml")

		// Check if either config file already exists
		if _, err := os.Stat(configPathYaml); err == nil {
			return fmt.Errorf("config file already exists at %s", configPathYaml)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("failed to check config file status at %s: %w", configPathYaml, err)
		}

		if _, err := os.Stat(configPathYml); err == nil {
			return fmt.Errorf("config file already exists at %s", configPathYml)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("failed to check config file status at %s: %w", configPathYml, err)
		}

		// Create directory if it doesn't exist
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// Write the file
		if err := os.WriteFile(configPathYaml, []byte(defaultConfigTemplate), 0644); err != nil {
			return fmt.Errorf("failed to write config file: %w", err)
		}

		slog.Info(fmt.Sprintf("Default config file created at %s\n", configPathYaml))
		return nil
	},
}
