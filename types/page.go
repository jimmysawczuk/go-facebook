package types

type Page struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`

	Likes             int `json:"likes"`
	TalkingAboutCount int `json:"talking_about_count"`

	About       string `json:"about"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Link        string `json:"link"`
	Website     string `json:"website"`

	Cover struct {
		ID      string `json:"id"`
		Source  string `json:"source"`
		OffsetX int    `json:"offset_x"`
		OffsetY int    `json:"offset_y"`
	} `json:"cover"`

	CanPost     bool `json:"can_post"`
	IsPublished bool `json:"is_published"`
}
