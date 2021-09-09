package bg_test

import (
	"errors"
	. "github.com/hx/golib/bg"
	. "github.com/hx/golib/testing"
	"os"
	"testing"
)

type DummyService struct {
	state   string
	started chan struct{}
	stopped chan struct{}
	err     error
}

func NewDummyService() *DummyService {
	return &DummyService{
		state:   "not started",
		started: make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

func (s *DummyService) Run() error {
	s.state = "running"
	close(s.started)
	<-s.stopped
	return s.err
}

func (s *DummyService) Stop() {
	s.state = "stopped"
	close(s.stopped)
}

func (s *DummyService) Fail(err error) {
	s.err = err
	s.state = "failed"
	close(s.stopped)
}

func TestGroup_Wait(t *testing.T) {
	s1 := NewDummyService()
	s2 := NewDummyService()

	Equals(t, "not started", s1.state)
	Equals(t, "not started", s2.state)

	group := new(Group)
	group.Add(s1, s2)

	<-s1.started
	<-s2.started

	Equals(t, "running", s1.state)
	Equals(t, "running", s2.state)

	s1Err := errors.New("uhoh")
	s1.Fail(s1Err)

	groupErr := group.Wait()
	Assert(t, errors.Is(groupErr, s1Err), "group should return s1 error")
	if err, ok := groupErr.(*UnexpectedStop); ok {
		Equals(t, s1, err.Service)
		Equals(t, s1Err, err.Err)
	}

	Equals(t, "failed", s1.state)
	Equals(t, "stopped", s2.state)

	Equals(t, s1Err, s1.err)
	Equals(t, nil, s2.err)
}

func TestGroup_Stop(t *testing.T) {
	s1 := NewDummyService()
	s2 := NewDummyService()

	group := new(Group)
	group.Add(s1, s2)

	<-s1.started
	<-s2.started

	group.Stop()

	Equals(t, nil, group.Wait())

	Equals(t, "stopped", s1.state)
	Equals(t, "stopped", s2.state)

	Equals(t, nil, s1.err)
	Equals(t, nil, s2.err)
}

func TestGroup_StopOnSignal(t *testing.T) {
	s1 := NewDummyService()
	s2 := NewDummyService()

	group := new(Group)
	group.Add(s1, s2)

	<-s1.started
	<-s2.started

	group.StopOnSignal(os.Interrupt)

	p, err := os.FindProcess(os.Getpid())
	Ok(t, err)
	Ok(t, p.Signal(os.Interrupt))

	Equals(t, nil, group.Wait())

	Equals(t, "stopped", s1.state)
	Equals(t, "stopped", s2.state)

	Equals(t, nil, s1.err)
	Equals(t, nil, s2.err)
}

func TestGroup_Add(t *testing.T) {
	t.Run("after stop", func(t *testing.T) {
		s1 := NewDummyService()
		s2 := NewDummyService()

		group := new(Group)
		group.Add(s1)
		<-s1.started
		doa := errors.New("doa")
		s1.Fail(doa)
		group.Wait()
		group.Add(s2)
		err := group.Wait()

		Assert(t, errors.Is(err, doa), "error should be DOA")
		Equals(t, "failed", s1.state)
		Equals(t, "not started", s2.state)
	})
}
