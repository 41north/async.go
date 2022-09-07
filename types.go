// Package async provides constructs for various asynchronous patterns.
package async

// Future represents a value of type T that will be set at some time in the future.
type Future[T any] interface {
	// Get returns a response channel of size 1 for receiving the future value.
	// If the value has already been set it will already be available within the return channel.
	Get() <-chan T

	// Set sets the return value and notifies consumers. Consumers are notified once only,
	// with the return value indicating if Set was successful or not.
	Set(value T) bool
}

// Result is a simple wrapper for representing a value or an error.
type Result[T any] interface {
	// Unwrap deconstructs the contents of this Result into a tuple.
	Unwrap() (T, error)
}

// CountingSemaphore can be used to limit the amount of in-flight processes / tasks.
type CountingSemaphore interface {
	// Size returns the total number of tokens available withing this CountingSemaphore.
	Size() int32

	// Acquire attempts to acquire an amount of tokens from the semaphore, waiting until it is successful.
	Acquire(count int32)

	// TryAcquire attempts to acquire an amount of tokens from the semaphore and returns whether
	// it was successful or not.
	TryAcquire(count int32) bool

	// Release attempts to return a certain amount of tokens to the semaphore, waiting until it is successful.
	Release(count int32)

	// TryRelease attempts to return a certain amount of tokens to the semaphore and returns whether
	// it was successful or not.
	TryRelease(count int32) bool
}
