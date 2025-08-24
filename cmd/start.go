/*
Copyright © 2025 Takeru Furuse
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
	"github.com/tsuperis3112/pmdr/internal/client"
	"github.com/tsuperis3112/pmdr/internal/ipc"
)

// StartCmd represents the start command
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a new Pomodoro session",
	Long:  `Starts a new Pomodoro session. If the daemon is not running, it will be started.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if daemon is running
		_, err := client.Status()
		if err != nil {
			slog.Info("Daemon not running, starting it now...")
			// Assume error means daemon is not running. Attempt to start it.
			daemonArgs := []string{"daemon"}
			// Pass through persistent flags
			if cmd.Flags().Changed("config") {
				cfg, _ := cmd.Flags().GetString("config")
				daemonArgs = append(daemonArgs, "--config", cfg)
			}
			if cmd.Flags().Changed("log-level") {
				level, _ := cmd.Flags().GetString("log-level")
				daemonArgs = append(daemonArgs, "--log-level", level)
			}
			if cmd.Flags().Changed("log-path") {
				path, _ := cmd.Flags().GetString("log-path")
				daemonArgs = append(daemonArgs, "--log-path", path)
			}

			daemonCmd := exec.Command(os.Args[0], daemonArgs...)
			logFile, err := os.Create("/tmp/pmdr_daemon.log")
			if err != nil {
				return fmt.Errorf("failed to create daemon log file: %w", err)
			}
			defer func() {
				if err := logFile.Close(); err != nil {
					slog.Error("Failed to close daemon log file", "error", err)
				}
			}()
			daemonCmd.Stderr = logFile
			if err := daemonCmd.Start(); err != nil {
				return fmt.Errorf("failed to start daemon: %w", err)
			}
			// Give the daemon a moment to start up
			time.Sleep(500 * time.Millisecond)
			slog.Info("Daemon started.")
		}

		// Prepare args for the start command
		startArgs := &ipc.StartArgs{}

		if cmd.Flags().Changed("work") {
			val, _ := cmd.Flags().GetString("work")
			d, err := time.ParseDuration(val)
			if err != nil {
				return fmt.Errorf("invalid work duration: %w", err)
			}
			startArgs.WorkDuration = &d
		}
		if cmd.Flags().Changed("short-break") {
			val, _ := cmd.Flags().GetString("short-break")
			d, err := time.ParseDuration(val)
			if err != nil {
				return fmt.Errorf("invalid short-break duration: %w", err)
			}
			startArgs.ShortBreakDuration = &d
		}
		if cmd.Flags().Changed("long-break") {
			val, _ := cmd.Flags().GetString("long-break")
			d, err := time.ParseDuration(val)
			if err != nil {
				return fmt.Errorf("invalid long-break duration: %w", err)
			}
			startArgs.LongBreakDuration = &d
		}
		if cmd.Flags().Changed("cycles") {
			val, _ := cmd.Flags().GetInt("cycles")
			startArgs.PomoCycles = &val
		}

		if err := client.Start(startArgs); err != nil {
			return fmt.Errorf("failed to start session: %w", err)
		}
		slog.Info("Pomodoro session started.")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(StartCmd)

	// Flags for overriding config values
	StartCmd.Flags().StringP("work", "w", "", "Work session duration (e.g., 25m)")
	StartCmd.Flags().StringP("short-break", "s", "", "Short break duration (e.g., 5m)")
	StartCmd.Flags().StringP("long-break", "l", "", "Long break duration (e.g., 15m)")
	StartCmd.Flags().IntP("cycles", "c", 0, "Number of work cycles before a long break")
}
