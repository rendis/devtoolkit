package devtoolkit

import "sync"

func NewConcurrentWorkers(maxWorkers int) *ConcurrentWorkers {
	return &ConcurrentWorkers{
		maxWorkers: maxWorkers,
		ch:         make(chan struct{}, maxWorkers),
		closed:     false,
	}
}

type ConcurrentWorkers struct {
	maxWorkers int
	ch         chan struct{}
	wg         sync.WaitGroup
	err        error
	closed     bool
	closeOnce  sync.Once
}

func (cw *ConcurrentWorkers) Execute(fn func()) {
	if cw.closed {
		return
	}

	cw.ch <- struct{}{}
	cw.wg.Add(1)
	go func() {
		defer func() {
			<-cw.ch
			cw.wg.Done()
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
	cw.wg.Wait()
}

func (cw *ConcurrentWorkers) IsOpen() bool {
	return !cw.closed
}

func (cw *ConcurrentWorkers) GetError() error {
	return cw.err
}

func (cw *ConcurrentWorkers) close(err error) {
	cw.closeOnce.Do(func() {
		cw.closed = true
		cw.err = err
		close(cw.ch)
	})
}
