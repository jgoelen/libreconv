package lock

import (
	"context"
)

// Lock is a semaphore which can be acquired with a specific timeout
type Lock interface {
	Try(ctx context.Context) bool
	Unlock()
}

type lock struct {
	sem chan int
}

func (lck lock) Try(ctx context.Context) bool {
	select {
	case lck.sem <- 1:
		return true
	case <-ctx.Done():
		return false
	}
}
func (lck lock) Unlock() {
	<-lck.sem
}

// New creates a new Lock instance
func New() Lock {
	return &lock{sem: make(chan int, 1)}
}
