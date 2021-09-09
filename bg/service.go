package bg

// Service is any Task that can be signalled to stop.
type Service interface {
	Task

	// Stop should signal the service to shut down gracefully.
	Stop()
}
