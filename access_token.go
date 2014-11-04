package facebook

import (
	"fmt"
	"time"
)

type AccessToken struct {
	token       string
	valid       bool
	expires     time.Time
	permissions []string
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

// Figures out what permissions are attached to the current access token.
func (at *AccessToken) Lint() error {
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
