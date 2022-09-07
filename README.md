# go-async

![Build](https://github.com/41north/go-async/actions/workflows/ci.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/41north/go-async/badge.svg?branch=feat/readme)](https://coveralls.io/github/41north/go-async?branch=feat/readme)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Status: _EXPERIMENTAL_

This library is primarily intended as a directed learning exercise and eventual collection of utilities and patterns
for working asynchronously in Go.

## Documentation

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/41north/go-async)

Full `go doc` style documentation for the project can be viewed online without
installing this package by using the excellent GoDoc site here:
http://godoc.org/github.com/41north/go-async

You can also view the documentation locally once the package is installed with
the `godoc` tool by running `godoc -http=":6060"` and pointing your browser to
http://localhost:6060/pkg/github.com/41north/go-async

## Installation

```bash
$ go get -u github.com/41north/go-async
```

Add this import line to the file you're working in:

```Go
import "github.com/41north/go-async"
```

## Quick Start

### Future

A basic example:

```Go
// create a string future
f := NewFuture[string]()

// create a consumer channel
ch := f.Get()
go func() {
	println(fmt.Sprintf("Value: %s", <-ch))
}()

// set the value
f.Set("hello")
```

### Counting Semaphore

A basic example:

```go
// we create an input and output channel for work needing to be done
inCh := make(chan string, 128)
outCh := make(chan int, 128)

// we want a max of 10 in-flight processes
s := NewCountingSemaphore(10)

// we create more workers than tokens available
for i := 0; i < 100; i++ {
	go func() {
		for {
			// acquire a token, waiting until one is available
			s.Acquire(1)

			// consume from the input channel
			v, ok := <-inCh
			if !ok {
				// channel was closed
				return
			}

			// do some work and produce an output value
			outCh <- len(v)

			// you need to be careful about releasing, if possible perform it with defer
			s.Release(1)
		}
	}()
}

// generate some work and put it into the work queue
// ...
// ...
```

There are more examples available in the go doc.

## License

Go-async is licensed under the [Apache 2.0 License](LICENSE)

## Contact

If you want to get in touch drop us an email at [hello@41north.dev](mailto:hello@41north.dev)
