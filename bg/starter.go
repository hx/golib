package bg

// Starter is a task that can be started using Start, and then later joined using Wait.
type Starter interface {
	// Start starts the task in a new goroutine. It must not be called more than once.
	Start()

	// Wait blocks until the task finishes, and returns its error. Wait can be called before or after Start is called.
	// It can be called multiple times. If called after the task has already finished, it will return the task's error
	// immediately.
	Wait() error
}

type starter struct {
	task   Task
	signal chan struct{}
	err    error
}

func NewStarter(task Task) Starter {
	return &starter{
		task:   task,
		signal: make(chan struct{}),
	}
}

func (s *starter) Start() {
	go func() {
		s.err = s.task.Run()
		close(s.signal)
	}()
}

func (s *starter) Wait() error {
	<-s.signal
	return s.err
}
