# Logur adapter for [Logrus](https://github.com/sirupsen/logrus)

[![CircleCI](https://circleci.com/gh/logur/adapter-logrus.svg?style=svg)](https://circleci.com/gh/logur/adapter-logrus)
[![Coverage](https://gocover.io/_badge/logur.dev/adapter/logrus)](https://gocover.io/logur.dev/adapter/logrus)
[![Go Report Card](https://goreportcard.com/badge/logur.dev/adapter/logrus?style=flat-square)](https://goreportcard.com/report/logur.dev/adapter/logrus)
[![GolangCI](https://golangci.com/badges/github.com/logur/adapter-logrus.svg)](https://golangci.com/r/github.com/logur/adapter-logrus)
[![Go Version](https://img.shields.io/badge/go%20version-%3E=1.11-61CFDD.svg?style=flat-square)](https://github.com/logur/adapter-logrus)
[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/logur.dev/adapter/logrus)


## Installation

```bash
go get logur.dev/adapter/logrus
```


## Usage

```go
package main

import (
	"github.com/sirupsen/logrus"
	logrusadapter "logur.dev/adapter/logrus"
)

func main() {
	logger := logrusadapter.New(logrus.New())
}
```


## Development

When all coding and testing is done, please run the test suite:

``` bash
$ make check
```


## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
