package facebook

import (
	"github.com/jimmysawczuk/go-facebook/types"
)

// GetPage builds and executes a GraphRequest for /{page_identifier}, tries to marshal the response into a types.Page.
func (c *Client) GetPage(pageIdentifier string) (page types.Page, err error) {
	err = c.Get("/"+pageIdentifier, GraphQueryString{
		"fields": []string{"id,name,username,fan_count,talking_about_count,about,description,category,link,website,cover,can_post,is_published,likes"},
	}).SetVersion(Latest).Exec(&page)
	return page, err
}

// GetUser builds and executes a GraphRequest for /{page_identifier}, tries to marshal the response into a types.User.
func (c *Client) GetUser(userIdentifier string) (user types.User, err error) {
	err = c.Get("/"+userIdentifier, GraphQueryString{
		"fields": []string{"id,name,first_name,middle_name,last_name,gender,age_range,birthday,link,website,locale,timezone,cover,is_verified,verified,installed,updated_time"},
	}).SetVersion(Latest).Exec(&user)
	return user, err
}
