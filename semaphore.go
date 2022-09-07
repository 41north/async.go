package async

import (
	"runtime"
	"sync/atomic"
)

// NewCountingSemaphore creates a new semaphore with specified amount of available tokens.
func NewCountingSemaphore(size int32) CountingSemaphore {
	result := countingSemaphore{size: size}
	result.tokens.Store(size)
	return &result
}

type countingSemaphore struct {
	size   int32
	tokens atomic.Int32
}

func (l *countingSemaphore) Size() int32 {
	return l.size
}

func (l *countingSemaphore) Acquire(count int32) {
	for !l.TryAcquire(count) {
		runtime.Gosched()
	}
}

func (l *countingSemaphore) TryAcquire(count int32) bool {
	if l.tokens.Add(count*-1) < 0 {
		// acquire failed, cancel the attempt by adding back what was removed
		l.tokens.Add(count)
		return false
	}
	return true
}

func (l *countingSemaphore) Release(count int32) {
	for !l.TryRelease(count) {
		runtime.Gosched()
	}
}

func (l *countingSemaphore) TryRelease(count int32) bool {
	current := l.tokens.Load()
	update := current + count
	if update > l.size {
		return false
	}
	return l.tokens.CompareAndSwap(current, update)
}
