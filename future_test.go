package async_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/41north/async.go"

	"github.com/stretchr/testify/assert"
)

func ExampleFuture_basic() {
	// create a string future
	f := async.NewFuture[string]()

	// create a consumer channel
	ch := f.Get()
	go func() {
		println(fmt.Sprintf("Value: %s", <-ch))
	}()

	// set the value
	f.Set("hello")
}

func ExampleFuture_multiple() {
	// create some futures
	foo := async.NewFuture[string]()
	bar := async.NewFuture[string]()

	// compute in the background
	go func() {
		foo.Set("foo")
	}()

	go func() {
		foo.Set("bar")
	}()

	// wait for their results
	println(<-foo.Get())
	println(<-bar.Get())
}

func ExampleFuture_select() {
	// create some futures
	foo := async.NewFuture[string]()
	bar := async.NewFuture[string]()

	// compute their values in the background
	go func() {
		foo.Set("foo")
	}()

	go func() {
		bar.Set("bar")
	}()

	// create some consumer channels
	fooCh := foo.Get()
	barCh := bar.Get()

	// wait with timeout

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var result []string
	finished := false

	for {
		select {
		case <-ctx.Done():
			fmt.Println("timeout")
			finished = true
		case v, ok := <-fooCh:
			if ok {
				result = append(result, v)
			}
			finished = len(result) == 2
		case v, ok := <-barCh:
			if ok {
				result = append(result, v)
			}
			finished = len(result) == 2
		}

		if finished {
			// break out of the loop
			break
		}
	}

	// print all the results
	fmt.Println(result)
}

func ExampleNewFutureImmediate() {
	f := async.NewFutureImmediate("hello")
	println(<-f.Get())
}

func TestFuture_Set(t *testing.T) {
	f := async.NewFuture[string]()
	assert.True(t, f.Set("foo"))
	assert.False(t, f.Set("bar"))
}

func TestFuture_Get(t *testing.T) {
	expected := "hello"
	f := async.NewFuture[string]()

	var getters []<-chan string
	for i := 0; i < 1000; i++ {
		getters = append(getters, f.Get())
	}

	go func() {
		<-time.After(1 * time.Nanosecond)
		f.Set(expected)
	}()

	for _, getter := range getters {
		v, ok := <-getter
		assert.True(t, ok, "channel was closed")
		assert.Equal(t, expected, v)

		_, ok = <-getter
		assert.False(t, ok, "channel should have been closed")
	}
}

func TestFuture_GetAfterValueIsSet(t *testing.T) {
	f := async.NewFuture[string]()

	expected := "hello"
	f.Set(expected)

	actual := <-f.Get()
	assert.Equal(t, expected, actual)
}

func TestImmediateFuture_Get(t *testing.T) {
	expected := "foo"
	f := async.NewFutureImmediate[string](expected)
	for i := 0; i < 1000; i++ {
		assert.Equal(t, expected, <-f.Get())
	}
}

func TestImmediateFuture_Set(t *testing.T) {
	defer func() {
		assert.Equal(t, recover(), async.PanicSetOnImmediateFuture)
	}()
	f := async.NewFutureImmediate[string]("foo")
	f.Set("bar")
}

func BenchmarkFuture(b *testing.B) {
	count := 100

	futures := make([]async.Future[string], count)
	for i := 0; i < count; i++ {
		futures[i] = async.NewFuture[string]()
		go func(idx int) {
			<-time.After(500 * time.Millisecond)
			futures[idx].Set(fmt.Sprintf("Result_%d", idx))
		}(i)
	}

	for _, f := range futures {
		<-f.Get()
	}
}

func BenchmarkFuture_Get(b *testing.B) {
	f := async.NewFuture[string]()
	go func() {
		f.Set("hello")
	}()

	for i := 0; i < 100; i++ {
		<-f.Get()
	}
}
