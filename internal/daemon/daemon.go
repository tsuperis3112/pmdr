package daemon

import (
	"fmt"
	"log/slog"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/tsuperis3112/pmdr/internal/config"
	"github.com/tsuperis3112/pmdr/internal/ipc"
)

// Run starts the pmdr daemon.
func Run() error {
	slog.Info("Starting pmdr daemon")

	// Write PID file
	pidPath := ipc.GetPidPath()
	if err := os.WriteFile(pidPath, []byte(strconv.Itoa(os.Getpid())), 0644); err != nil {
		return fmt.Errorf("failed to write pid file: %w", err)
	}
	defer func() {
		if err := os.Remove(pidPath); err != nil {
			slog.Error("Failed to remove pid file", "error", err)
		}
	}()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	timer := NewTimer(cfg)
	service := NewPmdrService(timer)

	if err := rpc.RegisterName(ipc.ServiceName, service); err != nil {
		return err
	}

	socketPath := ipc.GetSocketPath()
	// Remove old socket file if it exists
	if err := os.RemoveAll(socketPath); err != nil {
		return err
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := listener.Close(); err != nil {
			slog.Error("Failed to close listener", "error", err)
		}
	}()

	// Create a ticker that will advance the timer state.
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Goroutine to handle the timer ticks.
	go func() {
		for range ticker.C {
			timer.Tick()
		}
	}()

	slog.Info("Daemon listening on", "socket", socketPath)

	// Handle signals for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		slog.Info("Shutting down daemon")
		if err := listener.Close(); err != nil {
			slog.Error("Failed to close listener", "error", err)
		}
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			// Check if the error is due to the listener being closed.
			// If so, we can exit gracefully.
			if _, ok := err.(*net.OpError); ok {
				slog.Info("Listener closed, daemon shutting down.")
				break
			}
			slog.Error("Failed to accept connection", "error", err)
			continue
		}
		go rpc.ServeConn(conn)
	}

	return nil
}
