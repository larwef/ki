package runner

import (
	"log"
	"sync"
)

// Runnable defines behaviour for objects used by Runner
type Runnable interface {
	Serve(signal chan bool)
	GracefulShutdown()
}

// Runner takes care of running the Runnables. The signal is used to stop all runnables and the wait group makes sure all
// runnables are able to shut down properly.
type Runner struct {
	runnables []Runnable
	signal    chan bool
	waitGroup sync.WaitGroup
}

// NewRunner returns a new Runner object
func NewRunner() *Runner {
	newRunner := &Runner{}
	newRunner.signal = make(chan bool)

	return newRunner
}

// Add adds a Runnable to be run when calling the Run function
func (r *Runner) Add(runnable Runnable) {
	r.runnables = append(r.runnables, runnable)
}

// Run starts all the Runnables and shuts them down on signal. Returns when all Runnables are shut down.
func (r *Runner) Run() {
	if len(r.runnables) == 0 {
		log.Println("No runnables registered. Stopping application...")
	}

	r.waitGroup.Add(len(r.runnables))
	for _, element := range r.runnables {
		go func(runnable Runnable) {
			runnable.Serve(r.signal)
			r.waitGroup.Done()
		}(element)
	}

	select {
	case <-r.signal:
		log.Println("Received signal. Preparing for shutdown.")
	}

	for _, runnable := range r.runnables {
		runnable.GracefulShutdown()
	}

	r.waitGroup.Wait()
}
