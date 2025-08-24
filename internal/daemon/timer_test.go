package daemon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tsuperis3112/pmdr/internal/config"
	"github.com/tsuperis3112/pmdr/internal/ipc"
)

// testTimer is a wrapper around the Timer that allows for time manipulation.
type testTimer struct {
	*Timer
	currentTime time.Time
}

func newTestTimer(cfg *config.Config) *testTimer {
	tt := &testTimer{
		currentTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	timer := NewTimer(cfg)
	timer.nowFunc = func() time.Time {
		return tt.currentTime
	}
	tt.Timer = timer
	return tt
}

func (tt *testTimer) advanceTime(d time.Duration) {
	tt.currentTime = tt.currentTime.Add(d)
	tt.Tick()
}

func TestTimerStateTransitions(t *testing.T) {
	baseConfig := &config.Config{
		WorkDuration:       10 * time.Second,
		ShortBreakDuration: 5 * time.Second,
		LongBreakDuration:  8 * time.Second,
		PomoCycles:         2,
	}

	t.Run("start timer from stopped state", func(t *testing.T) {
		tm := newTestTimer(baseConfig)
		tm.Start(&ipc.StartArgs{})

		status := tm.Status()
		assert.Equal(t, ipc.StateRunning, status.State)
		assert.Equal(t, ipc.TypeWork, status.SessionType)
		assert.Equal(t, 1, status.PomoCycle)
		assert.Equal(t, 10*time.Second, status.RemainingTime)
	})

	t.Run("work session completes and transitions to short break", func(t *testing.T) {
		tm := newTestTimer(baseConfig)
		tm.Start(&ipc.StartArgs{})

		tm.advanceTime(10 * time.Second)

		status := tm.Status()
		assert.Equal(t, ipc.StateRunning, status.State)
		assert.Equal(t, ipc.TypeShortBreak, status.SessionType)
		assert.Equal(t, 1, status.PomoCycle)
		assert.Equal(t, 5*time.Second, status.RemainingTime)
	})

	t.Run("short break completes and transitions to work", func(t *testing.T) {
		tm := newTestTimer(baseConfig)
		tm.Start(&ipc.StartArgs{})

		tm.advanceTime(10 * time.Second) // work
		tm.advanceTime(5 * time.Second)  // break

		status := tm.Status()
		assert.Equal(t, ipc.StateRunning, status.State)
		assert.Equal(t, ipc.TypeWork, status.SessionType)
		assert.Equal(t, 2, status.PomoCycle)
		assert.Equal(t, 10*time.Second, status.RemainingTime)
	})

	t.Run("final work session completes and transitions to long break", func(t *testing.T) {
		tm := newTestTimer(baseConfig)
		tm.Start(&ipc.StartArgs{})

		tm.advanceTime(10 * time.Second) // work 1
		tm.advanceTime(5 * time.Second)  // break 1
		tm.advanceTime(10 * time.Second) // work 2

		status := tm.Status()
		assert.Equal(t, ipc.StateRunning, status.State)
		assert.Equal(t, ipc.TypeLongBreak, status.SessionType)
		assert.Equal(t, 0, status.PomoCycle)
		assert.Equal(t, 8*time.Second, status.RemainingTime)
	})

	t.Run("pause and resume", func(t *testing.T) {
		tm := newTestTimer(baseConfig)
		tm.Start(&ipc.StartArgs{})

		tm.advanceTime(3 * time.Second)
		tm.Pause()
		assert.Equal(t, ipc.StatePaused, tm.Status().State)

		remainingBeforePause := tm.Status().RemainingTime

		// Advance time manually without ticking
		tm.currentTime = tm.currentTime.Add(5 * time.Second)

		// Remaining time should be the same because the timer is paused
		assert.Equal(t, remainingBeforePause, tm.Status().RemainingTime)

		tm.Resume()
		assert.Equal(t, ipc.StateRunning, tm.Status().State)

		// After resuming, the remaining time should still be the same initially
		assert.Equal(t, remainingBeforePause, tm.Status().RemainingTime)

		// Time should decrease again after resuming
		tm.advanceTime(1 * time.Second)
		assert.Equal(t, remainingBeforePause-time.Second, tm.Status().RemainingTime)
	})

	t.Run("stop timer", func(t *testing.T) {
		tm := newTestTimer(baseConfig)
		tm.Start(&ipc.StartArgs{})
		tm.advanceTime(3 * time.Second)
		tm.Stop()
		assert.Equal(t, ipc.StateStopped, tm.Status().State)
	})
}
