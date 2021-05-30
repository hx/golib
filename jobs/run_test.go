package jobs_test

import (
	"errors"
	. "github.com/hx/golib/jobs"
	. "github.com/hx/golib/testing"
	"sync"
	"testing"
)

func TestParallel(t *testing.T) {
	var (
		jobs = make([]Job, 3)
		wait sync.WaitGroup
		err  = errors.New("derp")
	)
	wait.Add(len(jobs))
	for i := range jobs {
		job := JobFunc(func() error {
			wait.Done()
			wait.Wait()
			return err
		})
		jobs[i] = &job

	}
	var events []Event
	for e := range Parallel(len(jobs), jobs...) {
		events = append(events, e)
	}
	t.Run("there should be 3 events per job", func(t *testing.T) {
		Equals(t, len(jobs)*3, len(events))
	})
	t.Run("the first group events should be EventQueued", func(t *testing.T) {
		for i, job := range jobs {
			event, ok := events[i].(*EventQueued)
			Assert(t, ok, "event should be of type *EventQueued")
			Equals(t, job, event.Job())
		}
	})
	t.Run("the second group of events should be EventStarted", func(t *testing.T) {
		for i := len(jobs); i < len(jobs)*2; i++ {
			_, ok := events[i].(*EventStarted)
			Assert(t, ok, "event should be of type *EventStarted")
		}
	})
	t.Run("the third group of events should be EventFinished with the returned error", func(t *testing.T) {
		for i := len(jobs) * 2; i < len(jobs)*3; i++ {
			event, ok := events[i].(*EventFinished)
			Assert(t, ok, "event should be of type *EventFinished")
			Equals(t, err, event.Error())
		}
	})
}
