package devtoolkit

import (
	"context"
	"errors"
	"reflect"
	"sync"
)

var ConcurrentFnsAlreadyRunning = errors.New("concurrent fns already running")
var ConcurrentFnsNilContext = errors.New("context must not be nil")
var ConcurrentFnsNilEmpty = errors.New("fns must not be nil or empty")

type ConcurrentFn func(ctx context.Context) (any, error)

// ConcurrentExec is a struct that allows to execute a slice of ConcurrentFn concurrently
type ConcurrentExec struct {
	running             bool
	results             []any
	errs                []error
	mtx                 sync.Mutex
	concurrencyWg       sync.WaitGroup
	concurrencyCtx      context.Context
	cancelConcurrencyFn context.CancelFunc
}

func (ce *ConcurrentExec) ExecuteFns(ctx context.Context, fns ...ConcurrentFn) (ConcurrentExecResponse, error) {
	if ctx == nil {
		return nil, ConcurrentFnsNilContext
	}

	if fns == nil || len(fns) == 0 {
		return nil, ConcurrentFnsNilEmpty
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
	if !val.IsNil() {
		ce.results[pos] = result
	}
}

func (ce *ConcurrentExec) blockExecution() error {
	ce.mtx.Lock()
	defer ce.mtx.Unlock()
	if ce.running {
		return ConcurrentFnsAlreadyRunning
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

// ConcurrentExecResponse is the interface returned by ExecuteFns to interact with the results of the concurrent execution
type ConcurrentExecResponse interface {
	Results() []any        // blocks until all fns are done
	Errors() []error       // blocks until all fns are done
	CancelExecution()      // cancels the execution of all fns
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
