package async_test

import (
	"errors"
	"testing"

	"github.com/41north/async.go"

	"github.com/stretchr/testify/assert"
)

func ExampleNewResult() {
	result := async.NewResultValue[string]("success")
	v, _ := result.Unwrap()
	println(v)
}

func ExampleNewResultErr() {
	result := async.NewResultErr[string](errors.New("failure"))
	_, err := result.Unwrap()
	panic(err)
}

func TestNewResult(t *testing.T) {
	r := async.NewResult[string]("hello", nil)
	value, err := r.Unwrap()
	assert.Equal(t, "hello", value)
	assert.Nil(t, err)

	expected := errors.New("something bad happened")
	r = async.NewResult[string]("", expected)
	value, err = r.Unwrap()
	assert.Equal(t, "", value)
	assert.Equal(t, expected, err)
}

func TestNewResultValue(t *testing.T) {
	r := async.NewResultValue[string]("hello")
	value, err := r.Unwrap()
	assert.Equal(t, "hello", value)
	assert.Nil(t, err)
}

func TestNewResultErr(t *testing.T) {
	expected := errors.New("something bad happened")
	r := async.NewResultErr[string](expected)
	value, err := r.Unwrap()
	assert.Equal(t, "", value)
	assert.Equal(t, expected, err)
}
