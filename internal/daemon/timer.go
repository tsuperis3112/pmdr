package daemon

import (
	"sync"
	"time"

	"github.com/tsuperis3112/pmdr/internal/config"
	"github.com/tsuperis3112/pmdr/internal/hook"
	"github.com/tsuperis3112/pmdr/internal/ipc"
	"github.com/tsuperis3112/pmdr/internal/sound"
)

// Timer is a state machine for the pomodoro timer.
// It is designed to be thread-safe and does not manage its own ticker.
type Timer struct {
	mu sync.Mutex

	globalConfig  *config.Config
	sessionConfig *config.Config // Overridden for the current session

	state            ipc.SessionState
	sessionType      ipc.SessionType
	startSessionTime time.Time
	nextSessionTime  time.Time
	pauseTime        time.Time // Time when the timer was paused
	pomoCycle        int

	nowFunc func() time.Time
}

// NewTimer creates a new Timer.
func NewTimer(cfg *config.Config) *Timer {
	return &Timer{
		globalConfig: cfg,
		state:        ipc.StateStopped,
		nowFunc:      time.Now,
	}
}

// Tick advances the timer by a given duration and handles state transitions.
func (t *Timer) Tick() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state != ipc.StateRunning {
		return
	}

	now := t.nowFunc()
	if !t.nextSessionTime.After(now) {
		t.handleSessionCompletion()
	}
}

// Status returns the current status of the timer.
func (t *Timer) Status() ipc.StatusReply {
	t.mu.Lock()
	defer t.mu.Unlock()

	var remainingTime time.Duration
	if t.state == ipc.StatePaused {
		remainingTime = t.nextSessionTime.Sub(t.pauseTime)
	} else {
		remainingTime = t.nextSessionTime.Sub(t.nowFunc())
	}

	if remainingTime < 0 {
		remainingTime = 0
	}

	return ipc.StatusReply{
		State:         t.state,
		SessionType:   t.sessionType,
		RemainingTime: remainingTime,
		EndTime:       t.nextSessionTime,
		PomoCycle:     t.pomoCycle,
	}
}

// Start begins a new session.
func (t *Timer) Start(args *ipc.StartArgs) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state == ipc.StateRunning {
		return
	}

	t.stopInternal()

	cfg := *t.globalConfig
	if args.WorkDuration != nil {
		cfg.WorkDuration = *args.WorkDuration
	}
	if args.ShortBreakDuration != nil {
		cfg.ShortBreakDuration = *args.ShortBreakDuration
	}
	if args.LongBreakDuration != nil {
		cfg.LongBreakDuration = *args.LongBreakDuration
	}
	if args.PomoCycles != nil {
		cfg.PomoCycles = *args.PomoCycles
	}
	t.sessionConfig = &cfg

	t.pomoCycle = 1
	t.startSession(ipc.TypeWork)
}

// Pause pauses the timer.
func (t *Timer) Pause() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state != ipc.StateRunning {
		return
	}
	t.state = ipc.StatePaused
	t.pauseTime = t.nowFunc()
}

// Resume resumes the timer.
func (t *Timer) Resume() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state != ipc.StatePaused {
		return
	}
	durationPaused := t.nowFunc().Sub(t.pauseTime)
	t.nextSessionTime = t.nextSessionTime.Add(durationPaused)
	t.state = ipc.StateRunning
}

// Stop stops the timer completely.
func (t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.stopInternal()
}

// stopInternal stops the timer without locking.
func (t *Timer) stopInternal() {
	t.state = ipc.StateStopped
	t.sessionConfig = nil
}

// startSession starts a new session of the given type.
func (t *Timer) startSession(st ipc.SessionType) {
	switch st {
	case ipc.TypeWork:
		sound.Notify(sound.Work)
	case ipc.TypeShortBreak:
		sound.Notify(sound.ShortBreak)
	case ipc.TypeLongBreak:
		sound.Notify(sound.LongBreak)
	}

	t.sessionType = st
	t.state = ipc.StateRunning
	now := t.nowFunc()

	t.startSessionTime = now
	switch st {
	case ipc.TypeWork:
		t.nextSessionTime = now.Add(t.sessionConfig.WorkDuration)
	case ipc.TypeShortBreak:
		t.nextSessionTime = now.Add(t.sessionConfig.ShortBreakDuration)
	case ipc.TypeLongBreak:
		t.nextSessionTime = now.Add(t.sessionConfig.LongBreakDuration)
	}
}

// handleSessionCompletion decides what to do after a session ends.
func (t *Timer) handleSessionCompletion() {
	if t.state == ipc.StateStopped {
		return
	}

	completedSession := t.sessionType

	switch completedSession {
	case ipc.TypeWork:
		go hook.Run(t.sessionConfig.Hooks.Work)
	case ipc.TypeShortBreak:
		go hook.Run(t.sessionConfig.Hooks.ShortBreak)
	case ipc.TypeLongBreak:
		go hook.Run(t.sessionConfig.Hooks.LongBreak)
	}

	if completedSession == ipc.TypeWork {
		if t.pomoCycle >= t.sessionConfig.PomoCycles {
			t.pomoCycle = 0
			t.startSession(ipc.TypeLongBreak)
		} else {
			t.startSession(ipc.TypeShortBreak)
		}
	} else {
		t.pomoCycle++
		t.startSession(ipc.TypeWork)
	}
}
