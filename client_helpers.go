package facebook

import (
	"github.com/jimmysawczuk/go-facebook/types"
)

// GetPage builds and executes a GraphRequest for /{page_identifier}, tries to marshal the response into a types.Page.
func (c *Client) GetPage(pageIdentifier string) (page types.Page, err error) {
	err = c.Get("/"+pageIdentifier, nil).Exec(&page)
	return page, err
}

// GetUser builds and executes a GraphRequest for /{page_identifier}, tries to marshal the response into a types.User.
func (c *Client) GetUser(userIdentifier string) (user types.User, err error) {
	err = c.Get("/"+userIdentifier, nil).Exec(&user)
	return user, err
}
