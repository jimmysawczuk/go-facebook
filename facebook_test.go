package facebook

import (
	"testing"
)

var fb *Client
var access_token = "CAAU9I96vsVkBAL4YjkLLGbA6I15WZBzqU075wAzriZCLFuNPMRhKkV1WXUW2o9HujYF05Xdyk2xC7b9gZCdLZAgmR3d8y1ZBpIxNECerpxc4e0fd63ZCN6wuT2dJndgoZBzj9dXZA6fuC1aoTWMAPi28ZBiOBgqKEZB2ipS6LUKX8IvRfXsXS8fWxAxc8efFm6mPgZD"

func init() {
	fb = New("1474599152759129", "40de603edd149f514312b632c15bfdd3")
	fb.SetAccessToken(access_token)
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

	result, err := fb.Get("/budweiser", nil)
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
