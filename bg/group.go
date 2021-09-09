package bg

import (
	"os"
	"os/signal"
	"sync"
)

// Group manages a group of services by running them in the background, and stopping them when Stop is called. If any
// service unexpectedly stops on its own, all other services will also be stopped.
type Group struct {
	wait     sync.WaitGroup
	stopped  bool
	error    error
	services []Service
	mutex    sync.Mutex
	signals  chan os.Signal
}

// Add starts the given services in new goroutines, and adds them to the group.
func (g *Group) Add(services ...Service) {
	g.wait.Add(len(services))
	g.services = append(g.services, services...)
	for _, service := range services {
		go g.run(service)
	}
}

// Wait returns when all services have stopped. If Stop was called, the first service to return a non-nil error from its
// Service.Run will have its error returned from Wait. If a service stopped unexpectedly, its error will be returned,
// wrapped in UnexpectedStop.
func (g *Group) Wait() error {
	g.wait.Wait()
	if g.error == expectedStop {
		return nil
	}
	return g.error
}

// Stop calls Service.Stop on all running services, starting with the most recently added service.
func (g *Group) Stop() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.signals != nil {
		signal.Stop(g.signals)
		g.signals = nil
	}

	if g.stopped {
		return
	}

	err, isErr := g.error.(*UnexpectedStop)
	if err == nil {
		g.error = expectedStop
	}

	for i := len(g.services) - 1; i >= 0; i-- {
		service := g.services[i]
		if isErr && err.Service == service {
			continue
		}
		service.Stop()
	}

	g.stopped = true
}

// StopOnSignal listens in a new goroutine for the given signals, and calls Stop if and when they are received. If Stop
// is called manually, the goroutine stops listening and terminates.
//
// Typically, you'll be using an os.Interrupt to catch CTRL+C:
//  g.StopOnSignal(os.Interrupt)
func (g *Group) StopOnSignal(sig ...os.Signal) {
	g.mutex.Lock()

	if g.signals != nil {
		signal.Stop(g.signals)
	}

	g.signals = make(chan os.Signal, len(sig))
	signal.Notify(g.signals, sig...)

	signals := g.signals

	g.mutex.Unlock()

	go func() {
		<-signals
		g.mutex.Lock()

		signal.Stop(g.signals)
		g.signals = nil

		g.mutex.Unlock()
		g.Stop()
	}()
}

func (g *Group) run(service Service) {
	err := service.Run()
	g.mutex.Lock()
	if g.error == nil {
		g.error = &UnexpectedStop{service, err}
		g.mutex.Unlock()
		g.Stop()
	} else {
		g.mutex.Unlock()
	}
	g.wait.Done()
}
