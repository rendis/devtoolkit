package devtoolkit

import (
	"context"
	"errors"
	"sync"
)

type ConcurrentFn func(ctx context.Context) (any, error)

var ConcurrentFnsAlreadyRunning = errors.New("concurrent fns already running")

type ConcurrentExec struct {
	mtx         sync.Mutex
	cancelCtxFn context.CancelFunc
	running     bool
	err         error
	results     []any
}

func (ce *ConcurrentExec) ExecuteFns(cxt context.Context, fns ...ConcurrentFn) ([]any, error) {
	var concurrentWg sync.WaitGroup

	if err := ce.executeFns(cxt, fns, &concurrentWg); err != nil {
		return nil, err
	}

	concurrentWg.Wait()
	return ce.results, ce.err
}

func (ce *ConcurrentExec) executeFns(ctx context.Context, fns []ConcurrentFn, concurrentWg *sync.WaitGroup) error {
	if err := ce.blockExecution(); err != nil {
		return err
	}

	ce.results = make([]any, len(fns))
	ctx, ce.cancelCtxFn = context.WithCancel(ctx)
	for i, fn := range fns {
		go ce.executorWorker(ctx, i, fn, concurrentWg)
	}
	return nil
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

func (ce *ConcurrentExec) executorWorker(ctx context.Context, pos int, fn ConcurrentFn, concurrentWg *sync.WaitGroup) {
	defer concurrentWg.Done()

	result, err := fn(ctx)
	if err != nil {
		ce.stopWithError(err)
		return
	}

	select {
	case <-ctx.Done():
		return
	default:
		ce.results[pos] = result
	}
}

func (ce *ConcurrentExec) stopWithError(err error) {
	ce.mtx.Lock()
	defer ce.mtx.Unlock()
	if ce.running {
		ce.cancelCtxFn()
		ce.err = err
		ce.running = false
	}
}
