// Package facebook implements a few functions that basically wrap Go's REST client to work with the Facebook Graph API.
package facebook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// The Facebook Client object.
type Client struct {
	appId       string
	secret      string
	accessToken AccessToken
}

// A specific error that's returned from Facebook if there's an error with a request to the Graph API.
type GraphError struct {
	Code       int    `json:"code"`
	Subcode    int    `json:"error_subcode"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	HTTPStatus int    `json:"http"`
}

type AccessToken struct {
	token       string
	valid       bool
	expires     time.Time
	permissions []string
}

type rawGraphError struct {
	Error GraphError `json:"error"`
}

type HTTPMethod string

const (
	Get  HTTPMethod = "GET"
	Post            = "POST"
	Put             = "PUT"
)

const graph_endpoint string = "https://graph.facebook.com"

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
	f.LintAccessToken()
}

// Gets the working access token.
func (f *Client) AccessToken() AccessToken {
	return f.accessToken
}

// Figures out what permissions are attached to the current access token.
func (f *Client) LintAccessToken() error {
	at := f.accessToken

	if at.token == "" {
		return fmt.Errorf("Access token not set")
	}

	result, err := f.api("/debug_token", Get, map[string]interface{}{
		"input_token":  at.token,
		"access_token": f.appId + "|" + f.secret,
	})

	if err == nil {
		if _, exists := result["data"]; exists {
			data := result["data"].(map[string]interface{})

			if valid, exists := data["is_valid"]; exists {
				at.valid = valid.(bool)
			}

			if expires, exists := data["expires_at"]; exists {
				at.expires = time.Unix(int64(expires.(float64)), 0)
			}

			if _, exists := data["scopes"]; exists {
				perms := data["scopes"].([]interface{})
				at.permissions = []string{}
				for _, v := range perms {
					at.permissions = append(at.permissions, fmt.Sprintf("%s", v))
				}
			}
		}

		f.accessToken = at
	}

	return err
}

// Makes a standard API call to the Graph API.
func (f *Client) api(url string, method HTTPMethod, params map[string]interface{}) (map[string]interface{}, error) {

	if params == nil {
		params = make(map[string]interface{})
	}

	if _, exists := params["access_token"]; !exists && !f.accessToken.Empty() && f.accessToken.Valid() {
		params["access_token"] = f.accessToken.token
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

	if resp.StatusCode != 200 {
		g_err, err := parseError(resp.StatusCode, buf)
		if err == nil {
			return nil, g_err
		} else {
			return nil, err
		}
	} else {
		result := make(map[string]interface{})
		err = json.Unmarshal(buf, &result)

		return result, err
	}
}

// A wrapper method for Api to make POST requests.
func (f *Client) Post(url string, params map[string]interface{}) (result map[string]interface{}, err error) {
	return f.api(url, Post, params)
}

// A wrapper method for Api to make GET requests.
func (f *Client) Get(url string, params map[string]interface{}) (result map[string]interface{}, err error) {
	return f.api(url, Get, params)
}

// A wrapper method for Api to make PUT requests.
func (f *Client) Put(url string, params map[string]interface{}) (result map[string]interface{}, err error) {
	return f.api(url, Put, params)
}

func getQueryString(params map[string]interface{}) string {
	values, _ := url.ParseQuery("")

	for key, value := range params {
		switch value.(type) {
		case string:
			values.Add(key, value.(string))
		case fmt.Stringer:
			values.Add(key, value.(fmt.Stringer).String())
		default:
			panic("Can't make a string!")
		}
	}

	result := values.Encode()

	if result != "" {
		result = "?" + result
	}

	return result
}

func NewGraphError(data map[string]interface{}) (fge GraphError, err error) {
	b, _ := json.Marshal(data)
	err = json.Unmarshal(b, &fge)

	return fge, err
}

func (e GraphError) String() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func (e GraphError) Error() string {
	return e.String()
}

func parseError(status_code int, buf []byte) (GraphError, error) {
	raw_error := rawGraphError{}
	err := json.Unmarshal(buf, &raw_error)

	if err == nil {
		return raw_error.Error, nil
	} else {
		return GraphError{}, err
	}
}

func (at AccessToken) String() string {
	return fmt.Sprintf("%s, valid: %t, expires: %s, perms: %s", at.token, at.valid, at.expires, at.permissions)
}

func (at AccessToken) Empty() bool {
	return at.token == ""
}

func (at AccessToken) Valid() bool {
	return at.valid
}
