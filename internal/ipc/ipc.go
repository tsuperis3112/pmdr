package ipc

import (
	"fmt"
	"net"
	"os"
	"time"
)

const (
	// ServiceName is the name of the RPC service.
	ServiceName = "PmdrService"
	// SocketName is the name of the socket file.
	SocketName = "pmdr.sock"
	// PidFileName is the name of the pid file.
	PidFileName = "pmdr.pid"
)

// SessionState represents the state of the timer.
// NOTE: It is better to use stringer
type SessionState int

const (
	StateRunning SessionState = iota
	StatePaused
	StateDone
	StateStopped // Idle
)

// SessionType represents the type of the session.
// NOTE: It is better to use stringer
type SessionType int

const (
	TypeWork SessionType = iota
	TypeShortBreak
	TypeLongBreak
)

// StartArgs holds the arguments for the Start RPC call.
// Pointers are used to distinguish between a zero value and a value that was not set.
type StartArgs struct {
	WorkDuration       *time.Duration
	ShortBreakDuration *time.Duration
	LongBreakDuration  *time.Duration
	PomoCycles         *int
}

// Args holds arguments for RPC calls that don't need any.
type Args struct{}

// StatusReply holds the response for the status RPC call.

// StatusReply holds the response for the status RPC call.
type StatusReply struct {
	State         SessionState
	SessionType   SessionType
	RemainingTime time.Duration
	EndTime       time.Time
	PomoCycle     int
}

func getRuntimePath(fileName string) string {
	if runtimeDir := os.Getenv("XDG_RUNTIME_DIR"); runtimeDir != "" {
		return fmt.Sprintf("%s/%s", runtimeDir, fileName)
	}
	return fmt.Sprintf("/tmp/%s", fileName)
}

// GetSocketPath returns the path to the socket file.
func GetSocketPath() string {
	return getRuntimePath(SocketName)
}

// GetPidPath returns the path to the pid file.
func GetPidPath() string {
	return getRuntimePath(PidFileName)
}

// Dial dials the daemon's RPC server.
func Dial() (net.Conn, error) {
	return net.Dial("unix", GetSocketPath())
}
