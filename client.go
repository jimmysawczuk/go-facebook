// Package facebook implements a few functions that basically wrap Go's REST client to work with the Facebook Graph API.
package facebook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	Version10                   = "1.0"
	Version20                   = "2.0"
	Version21                   = "2.1"
	Version22                   = "2.2"
)

const graph_endpoint string = "https://graph.facebook.com"

type GraphRequest struct {
	Method  HTTPMethod
	Version GraphAPIVersion

	Path  string
	Query GraphQueryString

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

func (r *GraphRequest) Exec(target interface{}) error {

	path := r.Path
	if r.Version != Unversioned {
		path = "/" + string(r.Version) + "/" + path
	}

	url := url.URL{
		Scheme:   "https",
		Host:     "graph.facebook.com",
		Path:     path,
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
			return fmt.Errorf("couldn't unmarshal response into Graph Error: %s\n", err, string(buf))
		}
	}

	if r.IsJSON {
		err = json.Unmarshal(buf, target)
		if err != nil {
			return fmt.Errorf("error unmarshaling response into %T: %s", target, err)
		}
	} else if _, ok := target.(*[]byte); ok {
		*(target.(*[]byte)) = buf
	}

	return nil
}

func (f *Client) Get(path string, params GraphQueryString) *GraphRequest {
	return f.NewGraphRequest(Get, path, params)
}

func (f *Client) Post(path string, params GraphQueryString) *GraphRequest {
	return f.NewGraphRequest(Post, path, params)
}

func (f *Client) Put(path string, params GraphQueryString) *GraphRequest {
	return f.NewGraphRequest(Put, path, params)
}

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

// // Makes a standard API call to the Graph API.
// func (f *Client) api(req_url string, method HTTPMethod, params map[string]interface{}) (map[string]interface{}, error) {

// 	if params == nil {
// 		params = make(map[string]interface{})
// 	}

// 	if _, exists := params["access_token"]; !exists && !f.accessToken.Empty() && f.accessToken.Valid() {
// 		params["access_token"] = f.accessToken.token
// 	}

// 	req_url = graph_endpoint + req_url + getQueryString(params)

// 	resp, err := http.Get(req_url)
// 	if err != nil {
// 		return nil, err
// 	}

// 	buf, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if resp.StatusCode != 200 {
// 		g_err, err := parseError(resp.StatusCode, buf)
// 		if err == nil {
// 			return nil, g_err
// 		} else {
// 			return nil, err
// 		}
// 	} else {

// 		map_result := make(map[string]interface{})
// 		err = json.Unmarshal(buf, &map_result)
// 		if err == nil {
// 			return map_result, err
// 		}

// 		values, err := url.ParseQuery(string(buf))
// 		if err == nil {
// 			for k, v := range values {
// 				if len(v) == 1 {
// 					map_result[k] = v[0]
// 				} else {
// 					map_result[k] = v
// 				}
// 			}

// 			return map_result, nil
// 		}

// 		return map_result, fmt.Errorf("couldn't parse response from graph")
// 	}
// }

// // A wrapper method for Api to make POST requests.
// func (f *Client) Post(url string, params map[string]interface{}) (result map[string]interface{}, err error) {
// 	return f.api(url, Post, params)
// }

// // A wrapper method for Api to make GET requests.
// func (f *Client) Get(url string, params map[string]interface{}) (result map[string]interface{}, err error) {
// 	return f.api(url, Get, params)
// }

// // A wrapper method for Api to make PUT requests.
// func (f *Client) Put(url string, params map[string]interface{}) (result map[string]interface{}, err error) {
// 	return f.api(url, Put, params)
// }
