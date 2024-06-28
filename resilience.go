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
	MaxRetries       int              // indicates the maximum number of retries. Default is 3.
	WaitTime         time.Duration    // indicates the wait time between retries. Default is 100ms.
	Backoff          bool             // indicates whether to use exponential backoff. Default is false.
	RawError         bool             // indicates whether to return the raw error or wrap it in a new error. Default is false.
	IsIgnorableError func(error) bool // indicates whether to ignore the error or not. Default is nil.
	ReturnIgnorable  bool             // indicates whether to return the ignorable error or not. Default is false.
}

// NewResilience returns a new Resilience instance with the provided options or defaults.
func NewResilience(options *ResilienceOptions) (Resilience, error) {
	if options == nil {
		options = &ResilienceOptions{}
	}

	if options.MaxRetries < 0 {
		return nil, errors.New("MaxRetries cannot be negative")
	}

	if options.MaxRetries == 0 {
		options.MaxRetries = defaultMaxRetries
	}

	if options.WaitTime < 0 {
		return nil, errors.New("WaitTime cannot be negative")
	}

	if options.WaitTime == 0 {
		options.WaitTime = defaultWaitTime
	}

	return &resilience{*options}, nil
}

type resilience struct {
	ResilienceOptions
}

func (r *resilience) RetryOperation(operation func() error) error {
	var lastErr error
	waitTime := r.WaitTime
	for i := 0; i < r.MaxRetries; i++ {
		lastErr = operation()
		if lastErr == nil {
			return nil
		}

		if r.IsIgnorableError != nil && r.IsIgnorableError(lastErr) {
			if r.ReturnIgnorable {
				return lastErr
			}
			return nil
		}

		if r.Backoff {
			time.Sleep(waitTime)
			waitTime *= 2 // exponential backoff.
		} else {
			time.Sleep(r.WaitTime)
		}
	}

	if r.RawError {
		return lastErr
	}
	return errors.Join(lastErr, errors.New(fmt.Sprintf("max retries exceeded (%d)", r.MaxRetries)))
}
