// Package facebook implements a few functions that basically wrap Go's REST client to work with the Facebook Graph API.
package facebook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

// Client is an object which represents a pathway to the Facebook Graph API.
type Client struct {
	appID       string
	secret      string
	accessToken AccessToken

	// The default Graph API version to use
	DefaultVersion GraphAPIVersion

	// The HTTP client to use when making API requests
	HTTPClient HTTPClient
}

// HTTPMethod is a string which represents an HTTP method (e.g. GET, POST or PUT)
type HTTPMethod string

// GraphAPIVersion is a string which represents the Graph API version to use (e.g. v2.4)
type GraphAPIVersion string

// HTTPClient is the interface which is used to communicate with Facebook. It defaults to http.DefaultClient, but is implemented this way to allow for
// middleware functionality later on.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Some useful constants for building requests
const (
	Get  HTTPMethod = "GET"
	Post            = "POST"
	Put             = "PUT"

	Unversioned GraphAPIVersion = ""
	Version10                   = "v1.0"
	Version20                   = "v2.0"
	Version21                   = "v2.1"
	Version22                   = "v2.2"
	Version23                   = "v2.3"
	Version24                   = "v2.4"
	Version25                   = "v2.5"
	Version26                   = "v2.6"
	Latest                      = "v2.6"
)

const graphEndpoint string = "https://graph.facebook.com"

// GraphRequest is an HTTP request to the Graph API
type GraphRequest struct {
	// The HTTP method used
	Method HTTPMethod

	// Defaults to the client's DefaultVersion.
	Version GraphAPIVersion

	// e.g. /me or /starbucks
	Path string

	// A url.Values representation of desired query string parameters, e.g. width=50&height=50
	Query GraphQueryString

	// True if the expected content-type of the return is application/json. If this is true, Exec() will try
	// to marshal the response as JSON into the target object. Otherwise, it will just set the target object
	// as a []byte.
	IsJSON bool

	client HTTPClient
}

// GraphQueryString is a query string for a GraphRequest
type GraphQueryString url.Values

// An empty Facebook API client with which you can make public requests or set an arbitrary access token.
var BlankClient = New("", "")

// New returns a new *Client. Pass empty strings here if you don't need the object to have your App ID or Secret.
func New(appID string, secret string) (f *Client) {
	f = new(Client)

	f.appID = appID
	f.secret = secret

	f.HTTPClient = http.DefaultClient

	f.DefaultVersion = Unversioned

	return f
}

// SetAccessToken sets the working access token.
func (f *Client) SetAccessToken(at string) {
	f.accessToken = AccessToken{token: at}
	f.accessToken.Lint(f)
}

// LintAccessToken is an alias for client.AccessToken().Lint().
func (f *Client) LintAccessToken() (err error) {
	return f.accessToken.Lint(f)
}

// AccessToken returns the working access token.
func (f *Client) AccessToken() AccessToken {
	return f.accessToken
}

// NewGraphRequest builds a new GraphRequest with the passed method, path and query string parameters. If no access token is passed,
// but one is set in the client, it will be appended automatically. Assumes the response will be application/json.
func (f *Client) NewGraphRequest(method HTTPMethod, path string, params GraphQueryString) *GraphRequest {
	if params == nil {
		params = GraphQueryString{}
	}

	r := GraphRequest{
		Path:    path,
		Method:  method,
		Query:   params,
		IsJSON:  true,
		Version: f.DefaultVersion,
		client:  f.HTTPClient,
	}

	if _, exists := params["access_token"]; !exists && f.accessToken.Empty() == false {
		url.Values(r.Query).Add("access_token", f.accessToken.token)
	}

	return &r
}

// SetVersion sets the Graph API version on the request.
func (r *GraphRequest) SetVersion(v GraphAPIVersion) *GraphRequest {
	r.Version = v
	return r
}

// Exec executes the given request. Returns a GraphError if the response from Facebook is an error, or just
// a normal error if something goes wrong before that or during unmarshaling.
func (r *GraphRequest) Exec(target interface{}) error {

	p := r.Path
	if r.Version != Unversioned {
		p = "/" + string(r.Version) + "/" + p
	}

	p = path.Clean(p)

	url := url.URL{
		Scheme:   "https",
		Host:     "graph.facebook.com",
		Path:     p,
		RawQuery: url.Values(r.Query).Encode(),
	}

	req, _ := http.NewRequest(string(r.Method), url.String(), nil)
	if r.IsJSON {
		req.Header.Add("Accept", "application/json")
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("error setting up http request: %s", err)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body")
	}

	if resp.StatusCode != 200 {
		errorTarget := GraphError{}
		err = json.Unmarshal(buf, &errorTarget)
		if err != nil {
			return fmt.Errorf("couldn't unmarshal response into Graph Error: %s\n\t%s", err, string(buf))
		}

		return errorTarget
	}

	if _, ok := target.(*[]byte); ok {
		*(target.(*[]byte)) = buf
		return nil
	} else if r.IsJSON {
		err = json.Unmarshal(buf, target)
		if err != nil {
			return fmt.Errorf("error unmarshaling response into %T: %s\n\n%s", target, err, string(buf))
		}
	} else {
		return fmt.Errorf("invalid target type for non-json response: %T", target)
	}

	return nil
}

// Get is a wrapper for client.NewGraphRequest(Get, path, params)
func (f *Client) Get(path string, params GraphQueryString) *GraphRequest {
	return f.NewGraphRequest(Get, path, params)
}

// Post is a wrapper for client.NewGraphRequest(Post, path, params)
func (f *Client) Post(path string, params GraphQueryString) *GraphRequest {
	return f.NewGraphRequest(Post, path, params)
}

// Put is a wrapper for client.NewGraphRequest(Put, path, params)
func (f *Client) Put(path string, params GraphQueryString) *GraphRequest {
	return f.NewGraphRequest(Put, path, params)
}

// GetAppAccessToken builds an app access token for the client ID/secret of the client.
func (f *Client) GetAppAccessToken() (string, error) {
	var err error

	req := f.Get("/oauth/access_token", GraphQueryString{
		"client_id":     []string{f.appID},
		"client_secret": []string{f.secret},
		"grant_type":    []string{"client_credentials"},
	})

	target_obj := struct {
		Token string `json:"access_token"`
		Type  string `json:"token_type"`
	}{}
	err = req.Exec(&target_obj)
	if err == nil {
		return target_obj.Token, nil
	}

	fmt.Println(err)

	target_raw := []byte{}
	err = req.Exec(&target_raw)
	if err != nil {
		return "", fmt.Errorf("invalid response")
	}

	vals, _ := url.ParseQuery(string(target_raw))
	str, exists := vals["access_token"]
	if !exists || str[0] == "" {
		return "", fmt.Errorf("access token wasn't in response")
	}
	return str[0], nil
}
