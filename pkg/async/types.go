package async

type Future[T any] interface {
	Get() <-chan T
	Set(value T) bool
}

type Result[T any] interface {
	Unwrap() (T, error)
}
