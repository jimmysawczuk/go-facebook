package fql

import (
	"fmt"
	"testing"
)

func init() {
	_ = fmt.Sprintf
}

func TestBasicQuery(t *testing.T) {
	fql := NewFQLQuery("SELECT uid, name FROM user WHERE uid IN (15504121, 774070614)")
	err := fql.Exec()
	if err != nil {
		t.Fail()
	}

	if len(fql.Result) != 2 {
		t.Fail()
	}
}

func TestParameterizedQuery(t *testing.T) {
	fql := NewFQLQuery("SELECT uid, name FROM user WHERE uid IN (%d, %d)", 15504121, 774070614)
	err := fql.Exec()
	if err != nil {
		t.Fail()
	}

	if len(fql.Result) != 2 {
		t.Fail()
	}
}
