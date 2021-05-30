package jobs

import "context"

// Context is passed as the argument to each job's Run function. It implements context.Context, and can be used
// to cancel jobs where possible.
//
// It also allows jobs to spawn sub-jobs using its Run, RunContext, Parallel, and ParallelContext functions, which
// return their own Events busses. Events from sub-jobs are sent first to this bus, and then to the bus of the
// parent job.
type Context struct {
	context.Context
	events Events
	job    Job
	parent Job
}

// Progress can be used to generate additional EventProgressed events, with implementation-specific payloads.
func (c *Context) Progress(payload interface{}) {
	c.events <- &EventProgressed{
		job:     c.job,
		payload: payload,
	}
}

// Run runs the given jobs one after the other. If any job returns an error, subsequent jobs are not run.
func (c *Context) Run(jobs ...Job) (events Events) {
	return c.RunContext(c.Context, jobs...)
}

// RunContext is identical to Run, but accepts a context.Context.
func (c *Context) RunContext(ctx context.Context, jobs ...Job) (events Events) {
	return c.delegate(run(ctx, jobs, c.job, 1, true))
}

// Parallel runs the given jobs using a queue and the number of worker specified by concurrency.
// Errors do not affect continuation.
func (c *Context) Parallel(concurrency int, jobs ...Job) (events Events) {
	return c.ParallelContext(c.Context, concurrency, jobs...)
}

// ParallelContext is identical to Parallel, but accepts a context.Context.
func (c *Context) ParallelContext(ctx context.Context, concurrency int, jobs ...Job) (events Events) {
	return c.delegate(run(ctx, jobs, c.job, concurrency, false))
}

func (c *Context) delegate(from Events) (to Events) {
	to = make(Events)
	go func() {
		for event := range from {
			to <- event       // Dispatch to child bus
			c.events <- event // Dispatch to parent bus
		}
		close(to)
	}()
	return
}
