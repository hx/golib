package jobs

import (
	"context"
	"sync"
)

// Run runs the given jobs one after the other. If any job returns an error, subsequent jobs are not run.
func Run(jobs ...Job) Events {
	return RunContext(context.Background(), jobs...)
}

// RunContext is identical to Run, but accepts a context.Context.
func RunContext(ctx context.Context, jobs ...Job) Events {
	return run(ctx, jobs, nil, 1, true)
}

// Parallel runs the given jobs using a queue and the number of worker specified by concurrency.
// Errors do not affect continuation.
func Parallel(concurrency int, jobs ...Job) Events {
	return ParallelContext(context.Background(), concurrency, jobs...)
}

// ParallelContext is identical to Parallel, but accepts a context.Context.
func ParallelContext(ctx context.Context, concurrency int, jobs ...Job) Events {
	return run(ctx, jobs, nil, concurrency, false)
}

func run(ctx context.Context, jobs []Job, parent Job, concurrency int, haltOnError bool) (events Events) {
	if haltOnError && concurrency != 1 {
		panic("haltOnError must not be true unless concurrency is 1")
	}
	events = make(Events)
	go func() {
		for _, job := range jobs {
			events <- &EventQueued{job, parent}
		}
		r := &runner{ctx, jobs, parent, events}
		if concurrency == 1 {
			r.series(haltOnError)
		} else if concurrency > 1 && concurrency < len(jobs) {
			r.concurrently(concurrency)
		} else {
			r.parallel()
		}
		close(events)
	}()
	return
}

type runner struct {
	ctx    context.Context
	jobs   []Job
	parent Job
	events Events
}

func (r *runner) one(job Job) (err error) {
	r.events <- &EventStarted{job}
	err = job.Run(&Context{
		Context: r.ctx,
		events:  r.events,
		job:     job,
		parent:  r.parent,
	})
	r.events <- &EventFinished{job, err}
	return
}

func (r *runner) series(haltOnError bool) {
	for _, job := range r.jobs {
		if r.ctx.Err() != nil {
			return
		}
		if r.one(job) != nil && haltOnError {
			return
		}
	}
}

func (r *runner) parallel() {
	var wait sync.WaitGroup
	wait.Add(len(r.jobs))
	for i := range r.jobs {
		job := r.jobs[i]
		go func() {
			r.one(job)
			wait.Done()
		}()
	}
	wait.Wait()
}

func (r *runner) concurrently(concurrency int) {
	var (
		queue = make(chan Job)
		wait  sync.WaitGroup
	)
	wait.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			for job := range queue {
				r.one(job)
			}
			wait.Done()
		}()
	}
	for _, job := range r.jobs {
		if r.ctx.Err() != nil {
			break
		}
		queue <- job
	}
	close(queue)
	wait.Wait()
}
