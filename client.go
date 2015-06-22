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

// The Facebook Client object.
type Client struct {
	appId       string
	secret      string
	accessToken AccessToken
}

type HTTPMethod string
type GraphAPIVersion string

const (
	Get  HTTPMethod = "GET"
	Post            = "POST"
	Put             = "PUT"

	Unversioned GraphAPIVersion = ""
	Version10                   = "v1.0"
	Version20                   = "v2.0"
	Version21                   = "v2.1"
	Version22                   = "v2.2"
)

const graph_endpoint string = "https://graph.facebook.com"

type GraphRequest struct {
	Method HTTPMethod

	// Defaults to unversioned.
	Version GraphAPIVersion

	// e.g. /me or /starbucks
	Path string

	// A url.Values representation of desired query string parameters, e.g. width=50&height=50
	Query GraphQueryString

	// True if the expected content-type of the return is application/json. If this is true, Exec() will try
	// to marshal the response as JSON into the target object. Otherwise, it will just set the target object
	// as a []byte.
	IsJSON bool
}

type GraphQueryString url.Values

// An empty Facebook API client with which you can make public requests or set an arbitrary access token.
var BlankClient *Client = New("", "")

// Returns a new Client. Pass empty strings here if you don't need the object to have your App ID or Secret.
func New(appId string, secret string) (f *Client) {
	f = new(Client)

	f.appId = appId
	f.secret = secret

	return f
}

// Sets the working access token.
func (f *Client) SetAccessToken(at string) {
	f.accessToken = AccessToken{token: at}
	f.accessToken.Lint(f)
}

func (f *Client) LintAccessToken() (err error) {
	return f.accessToken.Lint(f)
}

// Gets the working access token.
func (f *Client) AccessToken() AccessToken {
	return f.accessToken
}

// Builds a new GraphRequest with the passed method, path and query string parameters. If no access token is passed,
// but one is set in the client, it will be appended automatically. Assumes the response will be application/json.
func (f *Client) NewGraphRequest(method HTTPMethod, path string, params GraphQueryString) *GraphRequest {
	if params == nil {
		params = GraphQueryString{}
	}

	r := GraphRequest{
		Path:   path,
		Method: method,
		Query:  params,
		IsJSON: true,
	}

	if _, exists := params["access_token"]; !exists && f.accessToken.Empty() == false {
		url.Values(r.Query).Add("access_token", f.accessToken.token)
	}

	return &r
}

// Executes the request. Returns a GraphError if the response from Facebook is an error, or just
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

	http_client := http.DefaultClient
	resp, err := http_client.Do(req)
	if err != nil {
		return fmt.Errorf("error setting up http request: %s", err)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body")
	}

	if resp.StatusCode != 200 {
		error_target := GraphError{}
		err = json.Unmarshal(buf, &error_target)
		if err == nil {
			return error_target
		} else {
			return fmt.Errorf("couldn't unmarshal response into Graph Error: %s\n\t%s", err, string(buf))
		}
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

// A wrapper for client.NewGraphRequest(Get, path, params)
func (f *Client) Get(path string, params GraphQueryString) *GraphRequest {
	return f.NewGraphRequest(Get, path, params)
}

// A wrapper for client.NewGraphRequest(Post, path, params)
func (f *Client) Post(path string, params GraphQueryString) *GraphRequest {
	return f.NewGraphRequest(Post, path, params)
}

// A wrapper for client.NewGraphRequest(Put, path, params)
func (f *Client) Put(path string, params GraphQueryString) *GraphRequest {
	return f.NewGraphRequest(Put, path, params)
}

// Returns an app access token for the client ID/secret of the client.
func (f *Client) GetAppAccessToken() (string, error) {
	req := f.Get("/oauth/access_token", GraphQueryString{
		"client_id":     []string{f.appId},
		"client_secret": []string{f.secret},
		"grant_type":    []string{"client_credentials"},
	})
	req.IsJSON = false

	target := []byte{}
	err := req.Exec(&target)
	if err == nil {
		vals, _ := url.ParseQuery(string(target))
		if str, exists := vals["access_token"]; exists && str[0] != "" {
			return str[0], nil
		} else {
			return "", fmt.Errorf("access token wasn't in response")
		}
	} else {
		return "", fmt.Errorf("error executing request for access token: %s", err)
	}
}
