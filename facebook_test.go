package facebook

import (
	"testing"
)

var fb *Client

func init() {
	fb = New("1474599152759129", "40de603edd149f514312b632c15bfdd3")

	token, err := fb.GetAppAccessToken()
	if err == nil {
		fb.SetAccessToken(token)
	}

}

func TestPermissions(t *testing.T) {
	err := fb.LintAccessToken()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestGraphAPI(t *testing.T) {

	result, err := fb.Get("/zuck", nil)
	if result == nil && err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}
}

func TestBlockedEndpoint(t *testing.T) {

	result, err := fb.Get("/1474599152759129", nil)
	if result == nil && err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

}

func TestInvalidCall(t *testing.T) {
	_, err := fb.Get("/1", nil)
	switch err.(type) {
	case GraphError:
	default:
		t.Errorf("Expected GraphError here, got %T", err)
		t.Fail()
	}
}
