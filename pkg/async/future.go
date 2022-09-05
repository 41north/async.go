package async

import (
	"context"
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

func NewImmediateFuture[T any](value T) Future[T] {
	return immediateFuture[T]{value: value}
}

type future[T any] struct {
	ctx context.Context

	value       atomic.Pointer[T]
	valueUpdate chan T

	listeners       atomic.Pointer[[]chan T]
	listenersUpdate chan interface{}
}

func NewFuture[T any](ctx context.Context) Future[T] {

	f := future[T]{
		ctx:             ctx,
		valueUpdate:     make(chan T),
		listenersUpdate: make(chan interface{}, 1),
	}
	f.listeners.Store(&[]chan T{}) // init the listeners to an empty list
	go f.publishValue()            // start the publishing routine
	return &f
}

func (f *future[T]) publishValue() {
	var value *T
	publishedIdx := -1

	publishToListeners := func() {
		listeners := *f.listeners.Load()
		for i := publishedIdx + 1; i < len(listeners); i++ {
			listeners[i] <- *value
			publishedIdx = i
		}
	}

	// keep watching for additions to the listeners slice and publish the value to any new additions
	for {
		select {

		case <-f.ctx.Done():
			return // stop trying to publish

		case v, ok := <-f.valueUpdate:
			if !ok {
				return
			}
			value = &v
			publishToListeners()

		case _, ok := <-f.listenersUpdate:
			if !ok || value == nil {
				continue
			}
			publishToListeners()
		}
	}
}

func (f *future[T]) Get() <-chan T {
	return f.addListener()
}

func (f *future[T]) addListener() <-chan T {
	ch := make(chan T, 1) // we do not want rendezvous here as it will block the publishing routine
	for {
		oldListenersPtr := f.listeners.Load()
		oldListeners := *oldListenersPtr
		newListeners := make([]chan T, len(oldListeners)+1)

		copy(newListeners, oldListeners)
		newListeners[len(oldListeners)] = ch

		if f.listeners.CompareAndSwap(oldListenersPtr, &newListeners) {
			f.listenersUpdate <- struct{}{}
			break
		}
	}
	return ch
}

func (f *future[T]) Set(value T) bool {
	ok := f.value.CompareAndSwap(nil, &value)
	if ok {
		// wake up the publishing routine
		f.valueUpdate <- value
	}
	return ok
}
