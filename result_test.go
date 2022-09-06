package async

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleNewResult() {
	result := NewResult[string]("success")
	v, _ := result.Unwrap()
	println(v)
}

func ExampleNewResultErr() {
	result := NewResultErr[string](errors.New("failure"))
	_, err := result.Unwrap()
	panic(err)
}

func TestNewResult(t *testing.T) {
	r := NewResult[string]("hello")
	value, err := r.Unwrap()
	assert.Equal(t, "hello", value)
	assert.Nil(t, err)
}

func TestNewResultErr(t *testing.T) {
	expected := errors.New("something bad happened")
	r := NewResultErr[string](expected)
	value, err := r.Unwrap()
	assert.Equal(t, "", value)
	assert.Equal(t, expected, err)
}
