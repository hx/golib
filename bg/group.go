package bg

import (
	"errors"
	"os"
	"os/signal"
	"sync"
)

// Group manages a group of services by running them in the background, and stopping them when Stop is called. If any
// service unexpectedly stops on its own, all other services will also be stopped.
type Group struct {
	// If set, OnShutdown will be called before the Group starts stopping its running services. If the shutdown is
	// caused by an error, that error will be passed as OnShutdown's argument.
	OnShutdown func(err error)

	wait     sync.WaitGroup
	topMutex sync.Mutex
	stopped  bool
	error    error
	services []Service
	endMutex sync.Mutex
	signals  chan os.Signal
}

// Add starts the given services in new goroutines, and adds them to the group.
func (g *Group) Add(services ...Service) {
	g.topMutex.Lock()
	if !g.stopped {
		g.wait.Add(len(services))
		g.services = append(g.services, services...)
		for _, service := range services {
			go g.run(service)
		}
	}
	g.topMutex.Unlock()
}

// Check calls callback with true if all servers are running as expected, or false if the group is shutting down, or has
// shut down.
//
// It is safe to call Add from within callback. Calls Stop and StopOnSignal will deadlock, as will calls to Wait if the
// group is still alive.
func (g *Group) Check(callback func(groupIsAlive bool)) {
	g.endMutex.Lock()
	defer g.endMutex.Unlock()
	callback(!g.stopped)
}

// Wait returns when all services have stopped. If Stop was called, the first service to return a non-nil error from its
// Service.Run will have its error returned from Wait. If a service stopped unexpectedly, its error will be returned,
// wrapped in UnexpectedStop.
func (g *Group) Wait() error {
	g.wait.Wait()
	return g.error
}

// Stop calls Service.Stop on all running services, starting with the most recently added service.
func (g *Group) Stop() {
	g.endMutex.Lock()
	defer g.endMutex.Unlock()

	if g.stopped {
		return
	}

	if g.signals != nil {
		signal.Stop(g.signals)
		g.signals = nil
	}

	g.topMutex.Lock()
	g.stopped = true
	g.topMutex.Unlock()

	err, isErr := g.error.(*UnexpectedStop)

	if g.OnShutdown != nil {
		g.OnShutdown(err)
	}

	for i := len(g.services) - 1; i >= 0; i-- {
		service := g.services[i]
		if isErr && err.Service == service {
			continue
		}
		service.Stop()
	}
}

// StopOnSignal listens in a new goroutine for the given signals, and calls Stop if and when they are received. If Stop
// is called manually, the goroutine stops listening and terminates.
//
// Typically, you'll be using an os.Interrupt to catch CTRL+C:
//  g.StopOnSignal(os.Interrupt)
func (g *Group) StopOnSignal(sig ...os.Signal) {
	g.endMutex.Lock()

	if g.signals != nil {
		signal.Stop(g.signals)
	}

	g.signals = make(chan os.Signal, len(sig))
	signal.Notify(g.signals, sig...)

	signals := g.signals

	g.endMutex.Unlock()

	go func() {
		<-signals
		g.endMutex.Lock()

		signal.Stop(g.signals)
		g.signals = nil

		g.endMutex.Unlock()
		g.Stop()
	}()
}

func (g *Group) run(service Service) {
	err := service.Run()
	g.endMutex.Lock()
	if !g.stopped {
		if err == nil {
			err = errors.New("no error")
		}
		g.error = &UnexpectedStop{service, err}
		g.endMutex.Unlock()
		g.Stop()
	} else {
		g.endMutex.Unlock()
	}
	g.wait.Done()
}
