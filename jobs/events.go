package jobs

// Events is a synchronous event bus for jobs and their sub-jobs.
type Events chan Event

// Wait consumes the entire event bus, and returns the last non-nil EventFinished.Error. Combine Wait with Run
// to run simple halt-on-error sequences. For example:
//
//  err := Run(firstJob, secondJob, thirdJob).Wait()
func (c Events) Wait() (err error) {
	for e := range c {
		if e, ok := e.(*EventFinished); ok {
			err = e.Error()
		}
	}
	return
}

// Event signals that a job has been queued, started, progressed, or finished.
type Event interface {
	// Job is the job for which the event has occurred.
	Job() Job
}

// EventQueued occurs when an event is queued by a call to Run, RunContext, Parallel, or ParallelContext. When
// multiple jobs are queued, EventQueued events occur in the same order.
type EventQueued struct {
	job    Job
	parent Job
}

// Job is the job that was queued.
func (e *EventQueued) Job() Job { return e.job }

// Parent is the job that queued Job, if it was queued using another job's Context. Otherwise, Parent is nil.
func (e *EventQueued) Parent() Job { return e.parent }

// EventStarted occurs when a Job starts.
type EventStarted struct {
	job Job
}

// Job is the job that started.
func (e *EventStarted) Job() Job { return e.job }

// EventProgressed occurs when a job progresses. It is triggered by a call from the job to its Context.Progress
// function.
type EventProgressed struct {
	job     Job
	payload interface{}
}

// Job is the job that has progressed.
func (e *EventProgressed) Job() Job { return e.job }

// Payload is the job-specific payload that contains information about the job's progress. It is the argument
// given by the job to its Context.Progress.
func (e *EventProgressed) Payload() interface{} { return e.payload }

// EventFinished occurs when a job finishes.
type EventFinished struct {
	job Job
	err error
}

// Job is the job that finished.
func (e *EventFinished) Job() Job { return e.job }

// Error is the error returned by the job.
func (e *EventFinished) Error() error { return e.err }
