package facebook

import (
	"testing"
)

func TestGraphAPI(t *testing.T) {
	fb := BlankAPIClient

	result, err := fb.Api("/zuck", Get, make(map[string]interface{}))

	if result == nil && err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}
}
