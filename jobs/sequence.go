package jobs

// Sequence is an implementation of Job that runs a sequence of jobs using its Context.Run.
type Sequence []Job

// Run implements Job.
func (s Sequence) Run(ctx *Context) (err error) {
	err = ctx.Run(s...).Wait()
	if err == nil {
		err = ctx.Err()
	}
	return
}
