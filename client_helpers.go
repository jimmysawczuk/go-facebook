package facebook

import (
	"github.com/jimmysawczuk/go-facebook/types"
)

func (this *Client) GetPage(page_identifier string) (page types.Page, err error) {
	err = this.Get("/"+page_identifier, nil).Exec(&page)
	return page, err
}

func (this *Client) GetUser(user_identifier string) (user types.User, err error) {
	err = this.Get("/"+user_identifier, nil).Exec(&user)
	return user, err
}
