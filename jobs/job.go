package jobs

// Job is a single unit of work.
type Job interface {
	Run(ctx *Context) (err error)
}

// JobFuncContext converts a simple function that accepts a Context to a Job implementation.
type JobFuncContext func(ctx *Context) (err error)

// Run implements Job.
func (j JobFuncContext) Run(ctx *Context) (err error) { return j(ctx) }

// JobFunc converts a simple function to a Job implementation.
func JobFunc(job func() (err error)) JobFuncContext { return func(*Context) error { return job() } }
