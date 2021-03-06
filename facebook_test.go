package facebook

import (
	"encoding/json"
	"fmt"
	"testing"
)

var fb *Client

func init() {
	fb = New("1474599152759129", "40de603edd149f514312b632c15bfdd3")
	fb.DefaultVersion = Version26

	token, err := fb.GetAppAccessToken()
	if err == nil {
		fb.SetAccessToken(token)
	} else {
		fmt.Printf("error setting access token: %s", err)
	}
}

func TestRawUnmarshal(t *testing.T) {
	req := fb.Get("/starbucks", nil)
	target := []byte{}
	err := req.Exec(&target)

	target2 := struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	}{}

	err = json.Unmarshal(target, &target2)
	if err != nil || target2.ID != "22092443056" {
		t.Errorf("Error: %s, ID: %s", err, target2.ID)
	}
}

func TestPermissions(t *testing.T) {
	err := fb.LintAccessToken()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGraphAPI(t *testing.T) {

	req := fb.Get("/starbucks", nil)

	target := struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	}{}

	err := req.Exec(&target)

	if err != nil || target.ID != "22092443056" {
		t.Errorf("Error: %s, ID: %s", err, target.ID)
	}
}

func TestInvalidGraphCall(t *testing.T) {
	req := fb.Get("/1", nil)

	target := struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	}{}

	err := req.Exec(&target)

	if graphError, ok := err.(GraphError); !ok {
		t.Errorf("call should have failed and returned a GraphError, but returned %T instead: %s", err, err)
	} else if ok && graphError.Code != 803 {
		t.Errorf("call should have failed with code 803, instead failed with %d", graphError.Code)
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

func TestGetVersionedPage(t *testing.T) {
	req := fb.Get("/starbucks", GraphQueryString{
		"fields": []string{"id,likes,name,username"},
	}).SetVersion(Version25)
	target := []byte{}
	err := req.Exec(&target)

	target2 := struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Likes    int64  `json:"likes"`
		Username string `json:"username"`
	}{}

	err = json.Unmarshal(target, &target2)
	if err != nil || target2.ID != "22092443056" {
		t.Errorf("Error: %s, ID: %s, raw: %s", err, target2.ID, string(target))
	}
}
