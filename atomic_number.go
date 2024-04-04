package devtoolkit

import (
	"golang.org/x/exp/constraints"
	"sync"
)

// Number is a type constraint that allows any type which is either an Integer or a Float.
type Number interface {
	constraints.Integer | constraints.Float
}

// AtomicNumber provides a concurrency-safe way to work with numeric values.
type AtomicNumber[T Number] struct {
	value T
	mu    sync.RWMutex
}

// Get returns the current value of the AtomicNumber.
func (a *AtomicNumber[T]) Get() T {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.value
}

// Set updates the value of the AtomicNumber.
func (a *AtomicNumber[T]) Set(value T) {
	a.mu.Lock()
	a.value = value
	a.mu.Unlock()
}

// GreaterThan returns true if the AtomicNumber's value is greater than the provided value.
func (a *AtomicNumber[T]) GreaterThan(value T) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.value > value
}

// LessThan returns true if the AtomicNumber's value is less than the provided value.
func (a *AtomicNumber[T]) LessThan(value T) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.value < value
}

// EqualTo returns true if the AtomicNumber's value is equal to the provided value.
func (a *AtomicNumber[T]) EqualTo(value T) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.value == value
}

// Increment increases the AtomicNumber's value by 1.
func (a *AtomicNumber[T]) Increment() {
	a.mu.Lock()
	a.value++
	a.mu.Unlock()
}

// IncrementBy increases the AtomicNumber's value by the provided amount.
func (a *AtomicNumber[T]) IncrementBy(n T) {
	a.mu.Lock()
	a.value += n
	a.mu.Unlock()
}

// IncrementAndGet increases the AtomicNumber's value by 1 and returns the new value.
func (a *AtomicNumber[T]) IncrementAndGet() T {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.value++
	return a.value
}

// IncrementByAndGet increases the AtomicNumber's value by the provided amount and returns the new value.
func (a *AtomicNumber[T]) IncrementByAndGet(n T) T {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.value += n
	return a.value
}

// IncrementIf increases the AtomicNumber's value by 1 if the provided condition is true.
// It returns true if the value was incremented.
func (a *AtomicNumber[T]) IncrementIf(cond func(T) bool) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if cond(a.value) {
		a.value++
		return true
	}
	return false
}

// IncrementByIf increases the AtomicNumber's value by the provided amount if the provided condition is true.
// It returns true if the value was incremented.
func (a *AtomicNumber[T]) IncrementByIf(n T, cond func(T) bool) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if cond(a.value) {
		a.value += n
		return true
	}
	return false
}

// Decrement decreases the AtomicNumber's value by 1.
func (a *AtomicNumber[T]) Decrement() {
	a.mu.Lock()
	a.value--
	a.mu.Unlock()
}

// DecrementBy decreases the AtomicNumber's value by the provided amount.
func (a *AtomicNumber[T]) DecrementBy(n T) {
	a.mu.Lock()
	a.value -= n
	a.mu.Unlock()
}

// DecrementAndGet decreases the AtomicNumber's value by 1 and returns the new value.
func (a *AtomicNumber[T]) DecrementAndGet() T {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.value--
	return a.value
}

// DecrementByAndGet decreases the AtomicNumber's value by the provided amount and returns the new value.
func (a *AtomicNumber[T]) DecrementByAndGet(n T) T {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.value -= n
	return a.value
}

// DecrementIf decreases the AtomicNumber's value by 1 if the provided condition is true.
// It returns true if the value was decremented.
func (a *AtomicNumber[T]) DecrementIf(cond func(T) bool) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if cond(a.value) {
		a.value--
		return true
	}
	return false
}

// DecrementByIf decreases the AtomicNumber's value by the provided amount if the provided condition is true.
// It returns true if the value was decremented.
func (a *AtomicNumber[T]) DecrementByIf(n T, cond func(T) bool) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if cond(a.value) {
		a.value -= n
		return true
	}
	return false
}
