package async

type result[T any] struct {
	value T
	err   error
}

func (r result[T]) Unwrap() (T, error) {
	return r.value, r.err
}

// NewResult creates a successful result.
func NewResult[T any](value T) Result[T] {
	return result[T]{value: value}
}

// NewResultErr creates a failed result.
func NewResultErr[T any](err error) Result[T] {
	return result[T]{err: err}
}
