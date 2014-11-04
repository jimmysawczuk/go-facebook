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
