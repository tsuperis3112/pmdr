package daemon

import (
	"github.com/tsuperis3112/pmdr/internal/ipc"
)

// PmdrService is the RPC service for pmdr.
type PmdrService struct {
	timer *Timer
}

// NewPmdrService creates a new PmdrService.
func NewPmdrService(t *Timer) *PmdrService {
	return &PmdrService{timer: t}
}

// Start starts the timer.
func (s *PmdrService) Start(args *ipc.StartArgs, reply *struct{}) error {
	s.timer.Start(args)
	return nil
}

// Pause pauses the timer.
func (s *PmdrService) Pause(args *ipc.Args, reply *struct{}) error {
	s.timer.Pause()
	return nil
}

// Resume resumes the timer.
func (s *PmdrService) Resume(args *ipc.Args, reply *struct{}) error {
	s.timer.Resume()
	return nil
}

// Stop stops the timer.
func (s *PmdrService) Stop(args *ipc.Args, reply *struct{}) error {
	s.timer.Stop()
	return nil
}

// Status returns the current status of the timer.
func (s *PmdrService) Status(args *ipc.Args, reply *ipc.StatusReply) error {
	*reply = s.timer.Status()
	return nil
}
