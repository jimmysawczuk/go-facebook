package types

import (
	"encoding/json"
	"fmt"
	"time"
)

// AgeRange is an object representing an age range for a given user (see: https://developers.facebook.com/docs/graph-api/reference/age-range/).
type AgeRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// Birthday is a time.Time with some custom JSON marshaling methods.
type Birthday time.Time

// Gender is a string that represents a Facebook user's gender.
type Gender string

// Two of the more common genders.
const (
	Male   Gender = "male"
	Female        = "female"
)

// CoverPhoto is an object that has information about a Facebook user's cover photo (see: https://developers.facebook.com/docs/graph-api/reference/cover-photo/)
type CoverPhoto struct {
	ID      string `json:"id"`
	Source  string `json:"source"`
	OffsetX int    `json:"offset_x"`
	OffsetY int    `json:"offset_y"`
}

// User is an object that represents a Facebook user as of v2.6 (see: https://developers.facebook.com/docs/graph-api/reference/user).
// Some of the more commonly used fields are included, but you may need additional permissions from the given user to get them all.
type User struct {
	ID string `json:"id"`

	Name       string `json:"name"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`

	Gender   Gender   `json:"gender"`
	AgeRange AgeRange `json:"age_range"`
	Birthday Birthday `json:"birthday"`
	Email    string   `json:"email"`

	Link     string     `json:"link"`
	Website  string     `json:"website"`
	Locale   string     `json:"locale"`
	Timezone int        `json:"timezone"`
	Cover    CoverPhoto `json:"cover"`

	IsVerified bool `json:"is_verified"`
	Verified   bool `json:"verified"`
	Installed  bool `json:"installed"`

	UpdatedTime time.Time `json:"updated_time"`
}

// MarshalJSON marshals a Birthday into a MM/DD/YYYY format.
func (bd Birthday) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(bd).Format("01/02/2006")), nil
}

// UnmarshalJSON unmarshals a MM/DD/YYYY string into a Birthday object.
func (bd *Birthday) UnmarshalJSON(in []byte) (err error) {
	bdStr := ""
	err = json.Unmarshal(in, &bdStr)
	if err != nil {
		return fmt.Errorf("Error parsing birthday from string: %s", err)
	}
	t, err := time.Parse("01/02/2006", bdStr)
	if err != nil {
		return fmt.Errorf("Error parsing birthday into time.Time: %s", err)
	}

	*bd = Birthday(t)
	return nil
}
