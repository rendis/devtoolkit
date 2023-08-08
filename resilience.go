package devtoolkit

import (
	"errors"
	"fmt"
	"time"
)

var (
	defaultMaxRetries int = 3
	defaultWaitTime       = 100 * time.Millisecond
)

// Resilience provides an interface for retrying operations in case of failure.
type Resilience interface {
	RetryOperation(operation func() error) error
}

// ResilienceOptions contains configuration parameters for retry operations.
type ResilienceOptions struct {
	MaxRetries *int
	WaitTime   *time.Duration
	Backoff    bool // indicates whether to use exponential backoff.
}

// NewResilience returns a new Resilience instance with the provided options or defaults.
func NewResilience(options *ResilienceOptions) (Resilience, error) {
	if options == nil {
		options = &ResilienceOptions{}
	}

	if options.MaxRetries == nil {
		options.MaxRetries = &defaultMaxRetries
	} else if *options.MaxRetries < 0 {
		return nil, errors.New("MaxRetries cannot be negative")
	}

	if options.WaitTime == nil {
		options.WaitTime = &defaultWaitTime
	} else if *options.WaitTime < 0 {
		return nil, errors.New("WaitTime cannot be negative")
	}

	return &resilience{
		maxRetries: *options.MaxRetries,
		waitTime:   *options.WaitTime,
		backoff:    options.Backoff,
	}, nil
}

type resilience struct {
	maxRetries int
	waitTime   time.Duration
	backoff    bool
}

func (r *resilience) RetryOperation(operation func() error) error {
	var lastErr error
	waitTime := r.waitTime
	for i := 0; i < r.maxRetries; i++ {
		lastErr = operation()
		if lastErr == nil {
			return nil
		}

		if r.backoff {
			time.Sleep(waitTime)
			waitTime *= 2 // exponential backoff.
		} else {
			time.Sleep(r.waitTime)
		}
	}

	return errors.Join(lastErr, fmt.Errorf("max retries exceeded (%d)", r.maxRetries))
}
