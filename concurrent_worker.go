package devtoolkit

import "sync"

func NewConcurrentWorkers(maxWorkers int) *ConcurrentWorkers {
	return &ConcurrentWorkers{
		maxWorkers: maxWorkers,
		ch:         make(chan struct{}, maxWorkers),
	}
}

type ConcurrentWorkers struct {
	maxWorkers int
	closed     bool
	err        error
	ch         chan struct{}
	wg         sync.WaitGroup
	closeOnce  sync.Once
	mu         sync.Mutex
}

func (cw *ConcurrentWorkers) Execute(fn func()) {
	cw.mu.Lock()
	if cw.closed {
		return
	}
	cw.ch <- struct{}{}
	cw.mu.Unlock()

	cw.wg.Add(1)
	go func() {
		defer func() {
			cw.wg.Done()
			<-cw.ch
		}()
		fn()
	}()
}

func (cw *ConcurrentWorkers) Wait() {
	cw.wg.Wait()
	cw.close(nil)
}

func (cw *ConcurrentWorkers) Stop(err error) {
	cw.close(err)
}

func (cw *ConcurrentWorkers) IsOpen() bool {
	return !cw.closed
}

func (cw *ConcurrentWorkers) GetError() error {
	return cw.err
}

func (cw *ConcurrentWorkers) close(err error) {
	cw.closeOnce.Do(func() {
		cw.mu.Lock()
		cw.err = err
		cw.closed = true
		close(cw.ch)
		cw.mu.Unlock()
	})
}
