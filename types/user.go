package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type AgeRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type Birthday time.Time

type Gender string

const (
	Male   Gender = "male"
	Female        = "female"
)

type User struct {
	ID string `json:"id"`

	Name       string `json:"name"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`

	Username string `json:"username"`

	Gender   Gender   `json:"gender"`
	AgeRange AgeRange `json:"age_range"`
	Birthday Birthday `json:"birthday"`
	Email    string   `json:"email"`

	Link     string `json:"link"`
	Website  string `json:"website"`
	Locale   string `json:"locale"`
	Timezone int    `json:"timezone"`
	Cover    struct {
		ID      string `json:"id"`
		Source  string `json:"source"`
		OffsetX int    `json:"offset_x"`
		OffsetY int    `json:"offset_y"`
	} `json:"cover"`

	IsVerified bool `json:"is_verified"`
	Verified   bool `json:"verified"`
	Installed  bool `json:"installed"`

	UpdatedTime time.Time `json:"updated_time"`
}

func (this Birthday) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(this).Format("01/02/2006")), nil
}

func (this *Birthday) UnmarshalJSON(in []byte) (err error) {
	birthday_str := ""
	err = json.Unmarshal(in, &birthday_str)
	if err != nil {
		return fmt.Errorf("Error parsing birthday from string: %s", err)
	}
	t, err := time.Parse("01/02/2006", birthday_str)
	if err != nil {
		return fmt.Errorf("Error parsing birthday into time.Time: %s", err)
	}

	*this = Birthday(t)
	return nil
}
