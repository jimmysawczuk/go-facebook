# `go-facebook`

A Facebook SDK for Go.

[ ![Codeship Status for jimmysawczuk/go-facebook](https://codeship.com/projects/10f96b80-90bd-0132-f540-66130ee6610f/status?branch=master)](https://codeship.com/projects/61581)

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

	The MIT License (MIT)
	Copyright (C) 2013-2015 by Jimmy Sawczuk

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in
	all copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE
