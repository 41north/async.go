package async

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFuture_Get(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	expected := "hello"
	f := NewFuture[string](ctx)

	var getters []<-chan string
	for i := 0; i < 1000; i++ {
		getters = append(getters, f.Get())
	}

	go func() {
		<-time.After(1 * time.Nanosecond)
		f.Set(expected)
	}()

	for _, getter := range getters {
		assert.Equal(t, expected, <-getter)
	}
}

func TestFuture_GetAfterValueIsSet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	f := NewFuture[string](ctx)

	expected := "hello"
	f.Set(expected)

	actual := <-f.Get()
	assert.Equal(t, expected, actual)
}

func TestImmediateFuture_Get(t *testing.T) {
	expected := "foo"
	f := NewImmediateFuture[string](expected)
	for i := 0; i < 1000; i++ {
		assert.Equal(t, expected, <-f.Get())
	}
}

func TestImmediateFuture_Set(t *testing.T) {
	defer func() {
		assert.Equal(t, recover(), PanicSetOnImmediateFuture)
	}()
	f := NewImmediateFuture[string]("foo")
	f.Set("bar")
}

func BenchmarkFuture(b *testing.B) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	count := 100

	futures := make([]Future[string], count)
	for i := 0; i < count; i++ {
		futures[i] = NewFuture[string](ctx)
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
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	f := NewFuture[string](ctx)
	go func() {
		f.Set("hello")
	}()

	for i := 0; i < 100; i++ {
		<-f.Get()
	}
}
