// Package devtoolkit provides a collection of utilities for Golang development.
package devtoolkit

import (
	"context"
	"errors"
	"reflect"
	"sync"
)

// Error values that can be returned by ConcurrentExec.
var (
	ConcurrentExecAlreadyRunningErr = errors.New("concurrent fns already running")
	ConcurrentExecNilContextErr     = errors.New("context must not be nil")
	ConcurrentExecFnsNilOrEmptyErr  = errors.New("fns must not be nil or empty")
)

// ConcurrentFn represents a function that can be executed concurrently. The function receives a context
// and returns a result and an error.
type ConcurrentFn func(ctx context.Context) (any, error)

// ConcurrentExec allows to execute a slice of ConcurrentFn concurrently.
// The running state, results, errors and context for the concurrent execution are stored within the struct.
type ConcurrentExec struct {
	running             bool
	results             []any
	errs                []error
	mtx                 sync.Mutex
	concurrencyWg       sync.WaitGroup
	concurrencyCtx      context.Context
	cancelConcurrencyFn context.CancelFunc
}

func NewConcurrentExec() *ConcurrentExec {
	return &ConcurrentExec{}
}

// ExecuteFns receives a context and a slice of functions to execute concurrently.
// It returns a ConcurrentExecResponse interface and an error if execution could not be started.
func (ce *ConcurrentExec) ExecuteFns(ctx context.Context, fns ...ConcurrentFn) (ConcurrentExecResponse, error) {
	if ctx == nil {
		return nil, ConcurrentExecNilContextErr
	}

	if fns == nil || len(fns) == 0 {
		return nil, ConcurrentExecFnsNilOrEmptyErr
	}

	if err := ce.executeFns(ctx, fns); err != nil {
		return nil, err
	}
	return ce, nil
}

func (ce *ConcurrentExec) executeFns(ctx context.Context, fns []ConcurrentFn) error {
	if err := ce.blockExecution(); err != nil {
		return err
	}

	ce.init(ctx, len(fns))

	for i, fn := range fns {
		ce.concurrencyWg.Add(1)
		go ce.executorWorker(i, fn)
	}

	return nil
}

func (ce *ConcurrentExec) init(ctx context.Context, totalFns int) {
	ce.errs = make([]error, totalFns)
	ce.results = make([]any, totalFns)
	ce.concurrencyCtx, ce.cancelConcurrencyFn = context.WithCancel(ctx)
}

func (ce *ConcurrentExec) executorWorker(pos int, fn ConcurrentFn) {
	defer ce.concurrencyWg.Done()
	result, err := fn(ce.concurrencyCtx)
	ce.errs[pos] = err
	val := reflect.ValueOf(result)

	// if result is not pointer
	if val.Kind() != reflect.Ptr {
		ce.results[pos] = result
	}

	// if result is pointer and not nil
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		ce.results[pos] = result
	}
}

func (ce *ConcurrentExec) blockExecution() error {
	ce.mtx.Lock()
	defer ce.mtx.Unlock()
	if ce.running {
		return ConcurrentExecAlreadyRunningErr
	}
	ce.running = true
	return nil
}

func (ce *ConcurrentExec) unblockExecution() {
	ce.mtx.Lock()
	defer ce.mtx.Unlock()
	ce.running = false
	ce.cancelConcurrencyFn()
}

// ConcurrentExecResponse is the interface returned by ExecuteFns to interact with the results of the concurrent execution.
type ConcurrentExecResponse interface {
	// Results blocks until all functions are done and returns the results.
	Results() []any // blocks until all fns are done

	// Errors blocks until all functions are done and returns any errors that occurred.
	Errors() []error // blocks until all fns are done

	// GetNotNilErrors blocks until all functions are done and returns any errors that occurred that are not nil.
	GetNotNilErrors() []error // blocks until all fns are done

	// CancelExecution cancels the execution of all functions.
	CancelExecution() // cancels the execution of all fns

	// Done returns a channel that is closed when all functions are done.
	Done() <-chan struct{} // returns a channel that is closed when all fns are done
}

func (ce *ConcurrentExec) Results() []any {
	ce.concurrencyWg.Wait()
	ce.unblockExecution()
	return ce.results
}

func (ce *ConcurrentExec) Errors() []error {
	ce.concurrencyWg.Wait()
	ce.unblockExecution()
	return ce.errs
}

func (ce *ConcurrentExec) GetNotNilErrors() []error {
	ce.concurrencyWg.Wait()
	ce.unblockExecution()
	var notNilErrors []error
	for _, err := range ce.errs {
		if err != nil {
			notNilErrors = append(notNilErrors, err)
		}
	}
	return notNilErrors
}

func (ce *ConcurrentExec) CancelExecution() {
	ce.unblockExecution()
	ce.concurrencyWg.Wait()
}

func (ce *ConcurrentExec) Done() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		ce.concurrencyWg.Wait()
		ce.unblockExecution()
		close(done)
	}()
	return done
}
