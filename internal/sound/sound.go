package sound

import (
	"log/slog"
	"os/exec"
	"runtime"

	"github.com/gen2brain/beeep"
)

// Type defines the type of notification.
type Type int

const (
	Work Type = iota
	ShortBreak
	LongBreak
)

// Notify speaks a message using the OS's native TTS engine.
// It falls back to a beep sound if the TTS engine is not available.
// The sound is played in a separate goroutine to avoid blocking.
func Notify(soundType Type) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Recovered from panic in sound.Notify", "panic", r)
			}
		}()

		var message string
		switch soundType {
		case Work:
			message = "Work session started."
		case ShortBreak:
			message = "Time for a short break."
		case LongBreak:
			message = "Time for a long break."
		default:
			playBeep()
			return
		}

		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("say", message)
		case "linux":
			cmd = exec.Command("spd-say", message)
		case "windows":
			cmd = exec.Command("PowerShell", "-Command", "Add-Type -AssemblyName System.Speech; (New-Object System.Speech.Synthesis.SpeechSynthesizer).Speak('"+message+"');")
		default:
			playBeep()
			return
		}

		if err := cmd.Run(); err != nil {
			slog.Error("Failed to run TTS command, falling back to beep", "os", runtime.GOOS, "error", err)
			playBeep()
		}
	}()
}

func playBeep() {
	if err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration); err != nil {
		slog.Error("Failed to play beep sound", "error", err)
	}
}
