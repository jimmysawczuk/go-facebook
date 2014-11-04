package facebook

import (
	"fmt"
	"testing"
)

var fb *Client

func init() {
	fb = New("1474599152759129", "40de603edd149f514312b632c15bfdd3")

	token, err := fb.GetAppAccessToken()
	if err == nil {
		fb.SetAccessToken(token)
	} else {
		fmt.Printf("%s", err)
	}

}

func TestPermissions(t *testing.T) {
	err := fb.LintAccessToken()
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}
}

func TestGraphAPI(t *testing.T) {

	req := fb.Get("/zuck", nil)

	target := struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:username"`
	}{}

	err := req.Exec(&target)

	if err != nil || target.ID != "4" {
		t.Errorf("Error: %s, ID: %d", err, target.ID)
	}
}

func TestInvalidGraphCall(t *testing.T) {

	req := fb.Get("/zuckkkkkkk", nil)

	target := struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:username"`
	}{}

	err := req.Exec(&target)

	if _, ok := err.(GraphError); !ok {
		t.Errorf("This call should have failed and marshalled into a GraphError, but didn't: %s", err)
	}
}

func TestBlockedEndpoint(t *testing.T) {

	req := fb.Get("/1474599152759129", nil)
	target := map[string]interface{}{}

	err := req.Exec(&target)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

}
