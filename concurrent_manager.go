package devtoolkit

import (
	"log"
	"sync"
	"time"
)

var releaseCondFn = func(value int) bool {
	return value > 0
}

// ConcurrentManager is a structure that manages a dynamic pool of workers.
// It can adjust the number of active workers based on the workload, within the provided minimum and maximum limits.
type ConcurrentManager struct {
	// Minimum number of workers.
	min int

	// Maximum number of workers.
	max int

	// Previous maximum number of workers, used for calculations.
	prevMax int

	// Current maximum number of workers, adjusted based on workload.
	currentMax int

	// An atomic counter tracking the number of allocated workers.
	allocated AtomicNumber[int]

	// Rate at which the number of workers is increased when needed.
	workerIncreaseRate float64

	// Time period after which the number of workers is potentially increased.
	timeIncreasePeriod time.Duration

	// Ensures some actions are only performed once.
	once sync.Once

	// A channel used to manage worker allocation.
	workers chan struct{}

	// Waits for all workers to finish before shutting down.
	wg sync.WaitGroup
}

// NewConcurrentManager creates a new instance of ConcurrentManager with specified parameters.
// It ensures that the provided parameters are within acceptable ranges and initializes the manager.
func NewConcurrentManager(minWorkers, maxWorkers int, workerIncreaseRate float64, timeIncreasePeriod time.Duration) *ConcurrentManager {
	if minWorkers == 0 {
		minWorkers = 1
	}

	if minWorkers > maxWorkers {
		maxWorkers = minWorkers
	}

	if workerIncreaseRate <= 1 {
		workerIncreaseRate = 1.5
	}

	if timeIncreasePeriod < 1 {
		timeIncreasePeriod = time.Second
	}

	var cw = &ConcurrentManager{
		min:                minWorkers,
		max:                maxWorkers,
		workerIncreaseRate: workerIncreaseRate,
		timeIncreasePeriod: timeIncreasePeriod,
	}

	cw.init()
	return cw
}

// Allocate requests a new worker to be allocated.
// It blocks if the maximum number of workers has been reached, until a worker is released.
func (c *ConcurrentManager) Allocate() {
	c.once.Do(func() {
		go c.tickToIncrease()
	})
	c.workers <- struct{}{}
	c.wg.Add(1)
	c.allocated.Increment()
}

// Release frees up a worker, making it available for future tasks.
// It only releases a worker if the release condition is met, ensuring resources are managed efficiently.
func (c *ConcurrentManager) Release() {
	if c.allocated.DecrementIf(releaseCondFn) {
		<-c.workers
		c.wg.Done()
	}
}

// Wait blocks until all workers have finished their tasks.
// It ensures that all resources are properly cleaned up before shutting down or reinitializing the manager.
func (c *ConcurrentManager) Wait() {
	log.Printf("waiting for workers to finish")
	c.wg.Wait()
	log.Printf("all workers finished")
	c.init()
}

// init initializes or resets the ConcurrentManager, setting up its internal structures and workers.
func (c *ConcurrentManager) init() {
	c.prevMax = c.min
	c.currentMax = c.min
	c.once = sync.Once{}
	if c.workers != nil {
		close(c.workers)
	}
	c.workers = make(chan struct{}, c.max)
	for range c.max - c.min {
		c.workers <- struct{}{}
	}
}

// calculateNewMax calculates and sets a new maximum number of workers based on the current workload and increase rate.
// It returns true if the maximum was adjusted, false if it has reached the predefined maximum limit.
func (c *ConcurrentManager) tickToIncrease() {
	ticker := time.NewTicker(c.timeIncreasePeriod)
	defer ticker.Stop()

	for range ticker.C {
		if !c.calculateNewMax() {
			log.Printf("workers already at max: %d", c.max)
			return
		}

		var delta = c.currentMax - c.prevMax
		log.Printf("increasing workers from %d to %d (delta: %d)", c.prevMax, c.currentMax, delta)

		for i := 0; i < delta; i++ {
			<-c.workers
		}

	}
}

// concurrentManagerCleanup is a cleanup function that is called when the ConcurrentManager is garbage collected.
// It ensures that all resources, particularly the worker channel, are properly released.
func (c *ConcurrentManager) calculateNewMax() bool {
	if c.currentMax == c.max {
		return false
	}

	c.prevMax = c.currentMax
	c.currentMax = int(float64(c.currentMax) * c.workerIncreaseRate)
	if c.currentMax > c.max {
		c.currentMax = c.max
	}
	return true
}
