package facebook

import (
	"testing"
)

var fb *Client
var access_token = "CAAU9I96vsVkBAJZCSSXZA32j1ZCZAiCngyveel3rRTYfVBXZCwcinHIrs7Bv5MINdgYHEG5mAUqPCbU6N6ATCRzcZAaBBpIwE5ajSNOQJTZBke5KAKGZBC1uaZAYZB7eKlcojCBGxJcmHc2B4m29vpH8ZCZCctlhgL8n1UuqgBZAzxJS09sgPOZAlEA0kiU29dRA0xlFgZD"

func init() {
	fb = New("1474599152759129", "40de603edd149f514312b632c15bfdd3")
	fb.SetAccessToken(access_token)
}

func TestPermissions(t *testing.T) {
	res := fb.LintAccessToken()
	if res != nil {
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

func TestInvalidCall(t *testing.T) {
	_, err := fb.Get("/1", nil)
	switch err.(type) {
	case GraphError:
	default:
		t.Errorf("Expected GraphError here, got %T", err)
		t.Fail()
	}
}
