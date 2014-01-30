// Package facebook implements a few functions that basically wrap Go's REST client to work with the Facebook Graph API.
package facebook

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// The Facebook Client object.
type Facebook struct {
	appId       string
	secret      string
	accessToken string
}

// A specific error that's returned from Facebook if there's an error with a request to the Graph API.
type FacebookGraphError struct {
	Code    int    `json:"code"`
	Subcode int    `json:"error_subcode"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type HTTPMethod string

const (
	Get  HTTPMethod = "GET"
	Post            = "POST"
	Put             = "PUT"
)

const graph_endpoint string = "https://graph.facebook.com"

// An empty Facebook API client with which you can make public requests or set an arbitrary access token.
var BlankAPIClient *Facebook = New("", "")

// Returns a new instance of the Facebook object. Pass empty strings here if you don't need the object to have your App ID or Secret.
func New(appId string, secret string) (f *Facebook) {
	f = new(Facebook)

	f.appId = appId
	f.secret = secret

	return f
}

// Sets the working access token.
func (f *Facebook) SetAccessToken(at string) {
	f.accessToken = at
}

// Gets the working access token.
func (f *Facebook) GetAccessToken() string {
	return f.accessToken
}

// Figures out what permissions are attached to the current access token.
func (f *Facebook) GetAccessTokenInfo() (permissions []interface{}, err error) {
	if f.accessToken == "" {
		return nil, errors.New("No new access token provided")
	}

	result, _ := f.api("/me/permissions", Get, nil)

	if result["data"] != nil {
		permissions = result["data"].([]interface{})
	} else if result["error"] != nil {
		permissions = nil
		err, _ = NewFacebookGraphError(result["error"].(map[string]interface{}))
	}

	return
}

// Makes a standard API call to the Graph API.
func (f *Facebook) api(url string, method HTTPMethod, params map[string]interface{}) (result map[string]interface{}, err error) {

	if params == nil {
		params = make(map[string]interface{})
	}

	if f.accessToken != "" {
		params["access_token"] = f.accessToken
	}

	url = graph_endpoint + url + getQueryString(params)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	result = make(map[string]interface{})
	err = json.Unmarshal(buf, &result)

	return result, err
}

// A wrapper method for Api to make POST requests.
func (f *Facebook) Post(url string, params map[string]interface{}) (result map[string]interface{}, err error) {
	return f.api(url, Post, params)
}

// A wrapper method for Api to make GET requests.
func (f *Facebook) Get(url string, params map[string]interface{}) (result map[string]interface{}, err error) {
	return f.api(url, Get, params)
}

// A wrapper method for Api to make PUT requests.
func (f *Facebook) Put(url string, params map[string]interface{}) (result map[string]interface{}, err error) {
	return f.api(url, Put, params)
}

func getQueryString(params map[string]interface{}) string {
	values, _ := url.ParseQuery("")

	for key, value := range params {
		switch value.(type) {
		default:
			panic("Can't make a string!")
		case string:
			values.Add(key, value.(string))
		case fmt.Stringer:
			values.Add(key, value.(fmt.Stringer).String())
		}
	}

	result := values.Encode()

	if result != "" {
		result = "?" + result
	}

	return result
}

func NewFacebookGraphError(data map[string]interface{}) (fge FacebookGraphError, err error) {
	b, _ := json.Marshal(data)
	err = json.Unmarshal(b, &fge)

	return fge, err
}

func (e FacebookGraphError) String() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func (e FacebookGraphError) Error() string {
	return e.String()
}
