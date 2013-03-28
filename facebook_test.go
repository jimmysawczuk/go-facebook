package facebook

import (
	"testing"
)

func TestGraphAPI(t *testing.T) {
	fb := BlankAPIClient

	result, err := fb.Get("/zuck", nil)

	if result == nil && err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}
}
