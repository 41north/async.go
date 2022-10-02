package async_test

import (
	"sync"
	"testing"

	"github.com/41north/async.go"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/btree"
)

func ExampleCountingSemaphore() {
	// we create an input and output channel for work needing to be done
	inCh := make(chan string, 128)
	outCh := make(chan int, 128)

	// we want a max of 10 in-flight processes
	s := async.NewCountingSemaphore(10)

	// we create more workers than tokens available
	for i := 0; i < 100; i++ {
		go func() {
			for {
				// acquire a token, waiting until one is available
				s.Acquire(1)

				// consume from the input channel
				v, ok := <-inCh
				if !ok {
					// channel was closed
					return
				}

				// do some work and produce an output value
				outCh <- len(v)

				// you need to be careful about releasing, if possible perform it with defer
				s.Release(1)
			}
		}()
	}

	// generate some work and put it into the work queue
	// ...
	// ...
}

func TestCountingSemaphore_TryAcquireAndRelease(t *testing.T) {
	s := async.NewCountingSemaphore(2)

	assert.Equal(t, int32(2), s.Size())

	// tokens are available
	assert.True(t, s.TryAcquire(1))
	assert.True(t, s.TryAcquire(1))

	// tokens have been exhausted
	assert.False(t, s.TryAcquire(1))
	assert.False(t, s.TryAcquire(10))

	// release
	assert.True(t, s.TryRelease(1))

	// acquire again
	assert.True(t, s.TryAcquire(1))

	// exhausted again
	assert.False(t, s.TryAcquire(1))

	// try to release more than the size
	assert.False(t, s.TryRelease(100))
}

func TestCountingSemaphore_Pipeline(t *testing.T) {
	workCount := 1000000
	workerCount := 100

	semaphore := async.NewCountingSemaphore(32)

	inCh := make(chan int, 1024)
	outCh := make(chan int, 1024)

	wg := sync.WaitGroup{}
	wg.Add(1)

	// process the output and verify we see all the work items produced
	go func() {
		results := btree.Set[int]{}
		for {
			v := <-outCh
			results.Insert(v)
			if results.Len() == workCount {
				break
			}
		}

		for i := 0; i < results.Len(); i++ {
			v, ok := results.GetAt(i)
			assert.True(t, ok)
			assert.Equal(t, i, v)
		}

		wg.Done()
	}()

	// generate workers and start processing
	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				semaphore.Acquire(1)
				v, ok := <-inCh
				if !ok {
					return // channel was closed
				}
				outCh <- v
				semaphore.Release(1)
			}
		}()
	}

	for i := 0; i < workCount; i++ {
		// generate work
		inCh <- i
	}
	close(inCh)

	wg.Wait()
}
