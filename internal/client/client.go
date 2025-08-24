package client

import (
	"fmt"
	"log/slog"
	"net/rpc"
	"os"
	"strconv"
	"syscall"

	"github.com/tsuperis3112/pmdr/internal/ipc"
)

func newClient() (*rpc.Client, error) {
	conn, err := ipc.Dial()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to daemon: %w. Is the daemon running?", err)
	}
	return rpc.NewClient(conn), nil
}

func call(serviceMethod string, args interface{}, reply interface{}) error {
	client, err := newClient()
	if err != nil {
		// If the daemon is not running, we don't need to return an error for stop command.
		if serviceMethod == ipc.ServiceName+".Stop" {
			return nil
		}
		return err
	}
	defer func() {
		if err := client.Close(); err != nil {
			slog.Error("Failed to close client connection", "error", err)
		}
	}()

	return client.Call(serviceMethod, args, reply)
}

func Start(args *ipc.StartArgs) error {
	return call(ipc.ServiceName+".Start", args, &struct{}{})
}

func Pause() error {
	return call(ipc.ServiceName+".Pause", &ipc.Args{}, &struct{}{})
}

func Resume() error {
	return call(ipc.ServiceName+".Resume", &ipc.Args{}, &struct{}{})
}

func Stop() error {
	// First, try to gracefully stop the timer via RPC.
	_ = call(ipc.ServiceName+".Stop", &ipc.Args{}, &struct{}{})

	// Then, read the PID file and send a SIGTERM signal.
	pidPath := ipc.GetPidPath()
	pidBytes, err := os.ReadFile(pidPath)
	if err != nil {
		if os.IsNotExist(err) {
			// PID file not found, daemon is likely not running.
			return nil
		}
		return fmt.Errorf("failed to read pid file: %w", err)
	}

	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		return fmt.Errorf("invalid pid in pid file: %w", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		// Process not found, it might have already exited.
		return nil
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		// On Unix, an error here often means the process is already gone.
		if err.Error() == "os: process already finished" {
			return nil
		}
		return fmt.Errorf("failed to send signal to daemon: %w", err)
	}

	return nil
}

func Status() (*ipc.StatusReply, error) {
	var reply ipc.StatusReply
	err := call(ipc.ServiceName+".Status", &ipc.Args{}, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}
