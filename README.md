# go-facebook

A basic Facebook SDK for Go.

[ ![travis-ci status for jimmysawczuk/go-facebook](https://travis-ci.org/jimmysawczuk/go-facebook.svg)](https://travis-ci.org/jimmysawczuk/go-facebook) [![GoDoc](https://godoc.org/github.com/jimmysawczuk/go-facebook?status.svg)](https://godoc.org/github.com/jimmysawczuk/go-facebook) [![Go Report Card](https://goreportcard.com/badge/github.com/jimmysawczuk/go-facebook)](https://goreportcard.com/report/github.com/jimmysawczuk/go-facebook)

## Installation

Install this package by typing `go get github.com/jimmysawczuk/go-facebook` in your terminal. You can then use it in your import statement like so:

```go
import (
	"github.com/jimmysawczuk/go-facebook"
)
```

## Example

```go
package main

import (
    "github.com/jimmysawczuk/go-facebook"
    "github.com/jimmysawczuk/go-facebook/types"
    "fmt"
)

func main() {
    fb := facebook.New("<app id>", "<secret>")
    fb.SetAccessToken("<token>")

    user, err := fb.GetUser("me")
    fmt.Println(user, err)

    resp := types.Page{}
    err = fb.Get("/starbucks", nil).Exec(&resp)
    fmt.Println(resp, err)
}
```

## Documentation

You can find the latest godoc output for this repository at [GoDoc.org](http://godoc.org/github.com/jimmysawczuk/go-facebook).

## License

go-facebook is released under [the MIT license](https://github.com/jimmysawczuk/go-facebook/blob/master/LICENSE).
