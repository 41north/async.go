package async

import (
	"sync/atomic"
)

const (
	PanicSetOnImmediateFuture = "you cannot set a value on an immediate future"
)

type immediateFuture[T any] struct {
	value T
}

func (f immediateFuture[T]) Get() <-chan T {
	ch := make(chan T, 1)
	ch <- f.value
	return ch
}

func (f immediateFuture[T]) Set(_ T) bool {
	panic(PanicSetOnImmediateFuture)
}

// NewFutureImmediate creates a future of type T that has a value that is already set.
func NewFutureImmediate[T any](value T) Future[T] {
	return immediateFuture[T]{value: value}
}

type future[T any] struct {
	value atomic.Pointer[T]

	consumers    atomic.Pointer[[]chan T]
	publishedIdx atomic.Int32
}

// NewFuture creates a new future of type T.
func NewFuture[T any]() Future[T] {
	f := future[T]{}
	f.consumers.Store(&[]chan T{}) // init the consumers to an empty list
	return &f
}

func (f *future[T]) Get() <-chan T {
	defer f.tryNotifyConsumers()
	return f.addConsumer()
}

func (f *future[T]) addConsumer() <-chan T {
	// we do not want rendezvous here as it will block the caller of tryNotifyConsumers()
	ch := make(chan T, 1)

	for {
		oldListenersPtr := f.consumers.Load()
		oldListeners := *oldListenersPtr
		newListeners := make([]chan T, len(oldListeners)+1)

		copy(newListeners, oldListeners)
		newListeners[len(oldListeners)] = ch

		if f.consumers.CompareAndSwap(oldListenersPtr, &newListeners) {
			break
		}
	}
	return ch
}

func (f *future[T]) tryNotifyConsumers() {
	value := f.value.Load()

	// check if the value has been set
	if value == nil {
		return
	}

	// determine if there are any new consumers that have not been notified
	consumers := *f.consumers.Load()
	consumerCount := len(consumers)

	publishedIdx := f.publishedIdx.Load()
	newConsumers := consumers[publishedIdx:consumerCount]

	if len(newConsumers) == 0 {
		// no new consumers
		return
	}

	if !f.publishedIdx.CompareAndSwap(publishedIdx, int32(consumerCount)) {
		// concurrency guard to prevent duplicate publications
		return
	}

	// notify the new consumers
	go func() {
		for _, consumer := range newConsumers {
			consumer <- *value
			close(consumer)
		}
	}()
}

func (f *future[T]) Set(value T) bool {
	ok := f.value.CompareAndSwap(nil, &value)
	if ok {
		f.tryNotifyConsumers()
	}
	return ok
}
