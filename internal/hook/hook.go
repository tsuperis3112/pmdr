package hook

import (
	"log/slog"
	"os/exec"
)

// Run executes the given commands in the background.
func Run(commands []string) {
	if len(commands) == 0 {
		return
	}

	for _, cmdStr := range commands {
		go func(c string) {
			cmd := exec.Command("sh", "-c", c)
			if err := cmd.Start(); err != nil {
				slog.Error("Failed to start hook command", "error", err, "command", c)
			}

			// We don't wait for the command to finish, but we should release resources.
			// The goroutine will exit after this, and the OS will handle the child process.
			_ = cmd.Wait() // This is to avoid zombie processes
			slog.Info("Executed hook command", "command", c, "pid", cmd.Process.Pid)
		}(cmdStr)
	}
}
