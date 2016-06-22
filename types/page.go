package types

// Page is an object that represents a Page object from the Graph API as of v2.6 (see https://developers.facebook.com/docs/graph-api/reference/page).
// Some commonly used fields are included.
type Page struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`

	FanCount          int `json:"fan_count"`
	TalkingAboutCount int `json:"talking_about_count"`

	About       string `json:"about"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Link        string `json:"link"`
	Website     string `json:"website"`

	Cover CoverPhoto `json:"cover"`

	CanPost     bool `json:"can_post"`
	IsPublished bool `json:"is_published"`

	Likes []Page `json:"data.likes"`
}
