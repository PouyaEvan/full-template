package resilience

import (
	"errors"
	"sync"
	"time"
)

// Simple Circuit Breaker implementation
type CircuitBreaker struct {
	mu           sync.RWMutex
	failureCount int
	threshold    int
	resetTimeout time.Duration
	lastFailure  time.Time
	state        string // "CLOSED", "OPEN", "HALF_OPEN"
}

func NewCircuitBreaker(threshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		threshold:    threshold,
		resetTimeout: resetTimeout,
		state:        "CLOSED",
	}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()
	if cb.state == "OPEN" {
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			cb.state = "HALF_OPEN"
		} else {
			cb.mu.Unlock()
			return errors.New("circuit breaker is open")
		}
	}
	cb.mu.Unlock()

	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailure = time.Now()
		if cb.failureCount >= cb.threshold {
			cb.state = "OPEN"
		}
		return err
	}

	// Success
	if cb.state == "HALF_OPEN" {
		cb.state = "CLOSED"
		cb.failureCount = 0
	} else if cb.state == "CLOSED" {
		cb.failureCount = 0
	}

	return nil
}
