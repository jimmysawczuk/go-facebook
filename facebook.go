// Package facebook implements a few functions that basically wrap Go's REST client to work with the Facebook Graph API.
package facebook

import (
	"encoding/json"
	"fmt"
	"time"
)

// The Facebook Client object.
type Client struct {
	appId       string
	secret      string
	accessToken AccessToken
}

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

func (f *Client) GetAppAccessToken() (string, error) {
	result, err := f.api("/oauth/access_token", Get, map[string]interface{}{
		"client_id":     f.appId,
		"client_secret": f.secret,
		"grant_type":    "client_credentials",
	})

	if access_token, set := result["access_token"]; set && err == nil {
		return access_token.(string), nil
	} else {
		return "", err
	}
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
