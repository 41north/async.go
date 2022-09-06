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

## Quick Start

### Future

Add this import line to the file you're working in:

```Go
import "github.com/41north/go-async"
```

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

There are more examples available in the go doc.

## License

Go-async is licensed under the [Apache 2.0 License](LICENSE)

## Contact

If you want to get in touch drop us an email at [hello@41north.dev](mailto:hello@41north.dev)
