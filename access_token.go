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
func (at *AccessToken) Lint(f *Client) error {
	if at.token == "" {
		return fmt.Errorf("Access token not set")
	}

	req := f.Get("/debug_token", GraphQueryString{
		"input_token":  []string{at.token},
		"access_token": []string{f.appId + "|" + f.secret},
	})

	target := struct {
		Data struct {
			AppId       string   `json:"app_id"`
			Valid       bool     `json:"is_valid"`
			Application string   `json:"application"`
			UserId      int64    `json:"user_id,string"`
			ExpiresAt   int64    `json:"expires_at"`
			Scopes      []string `json:"scopes"`
		} `json:"data"`
	}{}

	err := req.Exec(&target)
	if err == nil && target.Data.Valid {
		at.valid = target.Data.Valid
		at.expires = time.Unix(target.Data.ExpiresAt, 0)
		at.permissions = target.Data.Scopes
	}

	return err
}
