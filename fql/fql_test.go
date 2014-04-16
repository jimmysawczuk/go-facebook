package fql

import (
	"fmt"
	"testing"
)

func init() {
	_ = fmt.Sprintf
}

func TestBasicQuery(t *testing.T) {
	fql := NewQuery("SELECT uid, name FROM user WHERE uid IN (15504121, 774070614)")
	err := fql.Exec()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if len(fql.Result) != 2 {
		t.Fail()
	}
}

func TestParameterizedQuery(t *testing.T) {
	fql := NewQuery("SELECT uid, name FROM user WHERE uid IN (%d, %d)", 15504121, 774070614)
	err := fql.Exec()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if len(fql.Result) != 2 {
		t.Fail()
	}
}

func TestArrayQuery(t *testing.T) {
	fql := NewQuery("SELECT uid, name FROM user WHERE uid IN (%D)", []int{15504121, 774070614})
	err := fql.Exec()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if len(fql.Result) != 2 {
		t.Fail()
	}
}

func TestArrayCombinedQuery(t *testing.T) {
	fql := NewQuery("SELECT page_id, name, fan_count FROM page WHERE username IN (%S) AND fan_count > %d", []string{"coca-cola", "pepsi", "burgerking", "McDonalds"}, int64(3e7))
	err := fql.Exec()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if len(fql.Result) != 3 {
		t.Fail()
	}
}

func TestBadQuery(t *testing.T) {
	// invalid field
	fql := NewQuery("SELECT id, name, fan_count FROM page WHERE username IN ('pepsi') AND fan_count > %d", int64(3e7))
	err := fql.Exec()
	switch err.(type) {
	case Error:
	default:
		t.Errorf("Expected an error of type fql.Error, got: %#v", err)
		t.Fail()
	}
}
