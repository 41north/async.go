package async

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFuture_Get(t *testing.T) {
	expected := "hello"
	f := NewFuture[string]()

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
	f := NewFuture[string]()

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
	count := 100

	futures := make([]Future[string], count)
	for i := 0; i < count; i++ {
		futures[i] = NewFuture[string]()
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
	f := NewFuture[string]()
	go func() {
		f.Set("hello")
	}()

	for i := 0; i < 100; i++ {
		<-f.Get()
	}
}
