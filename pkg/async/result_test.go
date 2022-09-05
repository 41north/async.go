package async

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
