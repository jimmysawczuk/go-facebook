package facebook

import (
	"fmt"
	"time"
)

// AccessToken represents an Graph API access token
type AccessToken struct {
	token       string
	valid       bool
	expires     time.Time
	permissions []string
}

// String returns the access token as a string
func (at AccessToken) String() string {
	return at.token
}

// Info returns a formatted description of the access token, including the token itself, the permissions, the expiry and the validity.
func (at AccessToken) Info() string {
	return fmt.Sprintf("%s, valid: %t, expires: %s, perms: %s", at.token, at.valid, at.expires, at.permissions)
}

// Empty returns true if the access token is equal to the empty string.
func (at AccessToken) Empty() bool {
	return at.token == ""
}

// Valid returns the validity status for the token (but doesn't attempt to determine it if not present)
func (at AccessToken) Valid() bool {
	return at.valid
}

// Lint figures out what permissions are attached to the current access token.
func (at *AccessToken) Lint(f *Client) error {
	if at.token == "" {
		return fmt.Errorf("Access token not set")
	}

	req := f.Get("/debug_token", GraphQueryString{
		"input_token":  []string{at.token},
		"access_token": []string{f.appID + "|" + f.secret},
	})

	target := struct {
		Data struct {
			AppID       string   `json:"app_id"`
			Valid       bool     `json:"is_valid"`
			Application string   `json:"application"`
			UserID      int64    `json:"user_id,string"`
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
