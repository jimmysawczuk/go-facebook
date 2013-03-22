# `go-facebook`

A Facebook SDK for Go.

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
    "fmt"
)

func main() {
    fb := facebook.New("<app id>", "<secret>")
    fb.SetAccessToken("<token>")

    resp, err := fb.Api("/me", facebook.Get, nil)
    fmt.Println(resp, err)
}
```

## Documentation

Here's the output of `go doc` for now. I'll add better documentation soon.

```text
PACKAGE

package facebook
    import "go-facebook"

    Package facebook implements a few functions that basically wrap Go's
    REST client to work with the Facebook Graph API.

TYPES

type Facebook struct {
    // contains filtered or unexported fields
}
    The Facebook Client object.

var BlankAPIClient *Facebook = New("", "")
    An empty Facebook API client with which you can make public requests or
    set an arbitrary access token.

func New(appId string, secret string) (f *Facebook)
    Returns a new instance of the Facebook object. Pass empty strings here
    if you don't need the object to have your App ID or Secret.

func (f *Facebook) Api(url string, method HTTPMethod, params map[string]interface{}) (result map[string]interface{}, err error)
    Makes a standard API call to the Graph API.

func (f *Facebook) GetAccessToken() string
    Gets the working access token.

func (f *Facebook) GetAccessTokenInfo() (permissions []interface{}, err error)
    Figures out what permissions are attached to the current access token.

func (f *Facebook) SetAccessToken(at string)
    Sets the working access token.

type FacebookGraphError struct {
    Code    int
    Message string
    Type    string
}
    A specific error that's returned from Facebook if there's an error with
    a request to the Graph API.

func NewFacebookGraphError(code int, error_type string, message string) FacebookGraphError
    Instanciates a new Facebook Graph Error.

func (e FacebookGraphError) Error() string

func (e FacebookGraphError) String() string

type HTTPMethod string

const (
    Get  HTTPMethod = "GET"
    Post            = "POST"
    Put             = "PUT"
)
```

## License

	The MIT License (MIT)
	Copyright (C) 2013 by Jimmy Sawczuk

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
