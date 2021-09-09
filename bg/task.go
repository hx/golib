package bg

// Task is a process that returns an error if it fails.
type Task interface {
	// Run starts the task, and returns when it finishes.
	Run() error
}
