package display

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/tsuperis3112/pmdr/internal/ipc"
)

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	m := (d % time.Hour) / time.Minute
	s := (d % time.Minute) / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func formatSessionType(st ipc.SessionType) string {
	switch st {
	case ipc.TypeWork:
		return "Work"
	case ipc.TypeShortBreak:
		return "Short Break"
	case ipc.TypeLongBreak:
		return "Long Break"
	default:
		return "Unknown"
	}
}

func formatState(s ipc.SessionState) string {
	switch s {
	case ipc.StateRunning:
		return "Running"
	case ipc.StatePaused:
		return "Paused"
	case ipc.StateDone:
		return "Done"
	case ipc.StateStopped:
		return "Stopped"
	default:
		return "Unknown"
	}
}

// Status formats and prints the status reply
func Status(reply *ipc.StatusReply) {
	if reply.State == ipc.StateStopped {
		slog.Info("Timer is stopped.")
		return
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("[%s]", formatState(reply.State)))
	sb.WriteString(" ")
	sb.WriteString(formatSessionType(reply.SessionType))
	sb.WriteString(" ")
	sb.WriteString(formatDuration(reply.RemainingTime))

	if !reply.EndTime.IsZero() {
		sb.WriteString(fmt.Sprintf(" (ends at %s)", reply.EndTime.Format("15:04:05")))
	}

	if reply.SessionType == ipc.TypeWork {
		sb.WriteString(fmt.Sprintf(" (Cycle %d)", reply.PomoCycle))
	}

	slog.Info(sb.String())
}
