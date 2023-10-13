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
}

func (cw *ConcurrentWorkers) Close() {
	cw.closeOnce.Do(func() {
		cw.closed = true
		close(cw.ch)
	})
}

func (cw *ConcurrentWorkers) WaitAndClose() {
	cw.Wait()
	cw.Close()
}

func (cw *ConcurrentWorkers) Stop() {
	cw.Close()
}
