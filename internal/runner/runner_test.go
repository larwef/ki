package runner

import (
	"github.com/larwef/ki/test"
	"github.com/pkg/errors"
	"log"
	"testing"
)

// testRunnable is a mock object for testing Runner. mocks behaviour of blocking serve calls returning an error
type testRunnable struct {
	id       string
	shutdown chan bool
	err      error

	gracefulShutdownCalled bool
}

func newTestRunnable(id string, err error) *testRunnable {
	return &testRunnable{
		id:       id,
		shutdown: make(chan bool, 1),
		err:      err,
	}
}

func (tr *testRunnable) Serve(signal chan bool) {
	log.Printf("testRunnable %s serving\n", tr.id)
	if err := tr.serve(); err != nil {
		log.Printf("testRunnable %s exited serve with an error\n", tr.id)
		close(signal)
	}
}

func (tr *testRunnable) serve() error {
	select {
	case <-tr.shutdown:
		log.Printf("testRunnable %s received shutdown. Exiting Serve with error: %v\n", tr.id, tr.err)
		return tr.err
	}
}

func (tr *testRunnable) GracefulShutdown() {
	tr.shutdown <- true
	tr.gracefulShutdownCalled = true
}

// Tests happy path. Only one serve objects return error on shutdown.
func TestRunner_Run(t *testing.T) {
	runner := NewRunner()

	runnable1 := newTestRunnable("1", nil)
	runnable2 := newTestRunnable("2", nil)
	runnable3 := newTestRunnable("3", errors.New("Some Error"))
	runnable4 := newTestRunnable("4", nil)

	runner.Add(runnable1)
	runner.Add(runnable2)
	runner.Add(runnable3)
	runner.Add(runnable4)

	runnable3.shutdown <- true

	runner.Run()

	test.AssertEqual(t, runnable1.gracefulShutdownCalled, true)
	test.AssertEqual(t, runnable2.gracefulShutdownCalled, true)
	test.AssertEqual(t, runnable3.gracefulShutdownCalled, true)
	test.AssertEqual(t, runnable4.gracefulShutdownCalled, true)
}

// TODO: Uncoment when panic issue for serve is adressed. Se other todos
//// Test case where multiple serve objects returns error
//func TestRunner_Run_WithMultipleErrors(t *testing.T) {
//	runner := NewRunner()
//
//	runnable1 := newTestRunnable("1", nil)
//	runnable2 := newTestRunnable("2", nil)
//	runnable3 := newTestRunnable("3", errors.New("Some Error"))
//	runnable4 := newTestRunnable("4", errors.New("Some Error"))
//
//	runner.Add(runnable1)
//	runner.Add(runnable2)
//	runner.Add(runnable3)
//	runner.Add(runnable4)
//
//	runnable3.shutdown <- true
//
//	runner.Run()
//
//	test.AssertEqual(t, runnable1.gracefulShutdownCalled, true)
//	test.AssertEqual(t, runnable2.gracefulShutdownCalled, true)
//	test.AssertEqual(t, runnable3.gracefulShutdownCalled, true)
//	test.AssertEqual(t, runnable4.gracefulShutdownCalled, true)
//}
