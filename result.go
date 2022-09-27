package async

type result[T any] struct {
	value T
	err   error
}

func (r result[T]) Unwrap() (T, error) {
	return r.value, r.err
}

// NewResult creates a result instance with a provided value and error. It's sometimes more convenient to instantiate
// like this when implementing library code.
func NewResult[T any](value T, err error) Result[T] {
	return result[T]{value: value, err: err}
}

// NewResultValue creates a successful result.
func NewResultValue[T any](value T) Result[T] {
	return result[T]{value: value}
}

// NewResultErr creates a failed result.
func NewResultErr[T any](err error) Result[T] {
	return result[T]{err: err}
}
