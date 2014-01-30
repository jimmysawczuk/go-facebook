package facebook

import (
	"testing"
)

var fb Facebook = *BlankAPIClient

func TestGraphAPI(t *testing.T) {

	result, err := fb.Get("/zuck", nil)

	if result == nil && err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}
}

func TestInvalidCall(t *testing.T) {
	result, _ := fb.Get("/zuckkkkkkkkkkkk", nil)

	if _, exists := result["error"]; exists {
		_, err := NewFacebookGraphError(result["error"].(map[string]interface{}))
		if err != nil {
			t.Errorf("%s", err)
			t.Fail()
		}

	} else {
		t.Fail()
	}
}
