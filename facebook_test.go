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
	req := fb.Get("/1", nil)

	target := struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:username"`
	}{}

	err := req.Exec(&target)

	if graph_error, ok := err.(GraphError); !ok {
		t.Errorf("call should have failed and returned a GraphError, but returned %T instead: %s", err, err)
	} else if ok && graph_error.Code != 803 {
		t.Errorf("call should have failed with code 803, instead failed with %d", graph_error.Code)
	}
}

func TestBlockedEndpoint(t *testing.T) {

	req := fb.Get("/1474599152759129", nil)
	target := map[string]interface{}{}

	err := req.Exec(&target)
	if err != nil {
		t.Errorf("%s", err)
	}

}

func TestGetPage(t *testing.T) {
	pg, err := fb.GetPage("starbucks")
	if err != nil || pg.Username != "Starbucks" {
		t.Errorf("%s", err)
	}

	pg2, err := fb.GetPage("22092443056")
	if err != nil || pg2.Username != "Starbucks" {
		t.Errorf("%s", err)
	}

}

func TestGetUser(t *testing.T) {
	jimmy, err := fb.GetUser("15504121")
	if err != nil || jimmy.Name != "Jimmy Sawczuk" {
		t.Errorf("%s", err)
	}
}
