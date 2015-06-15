package fql

import (
	"fmt"
	"testing"
)

func init() {
	_ = fmt.Sprintf
}

func TestBasicQuery(t *testing.T) {
	fql := NewQuery("SELECT page_id, name FROM page WHERE username IN ('Starbucks', 'McDonalds')")
	err := fql.Exec()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if fql.Result.Len() != 2 {
		t.Fail()
	}
}

func TestParameterizedQuery(t *testing.T) {
	fql := NewQuery("SELECT page_id, name FROM page WHERE username IN (%s, %s)", "Starbucks", "McDonalds")
	err := fql.Exec()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if fql.Result.Len() != 2 {
		t.Fail()
	}
}

func TestArrayQuery(t *testing.T) {
	fql := NewQuery("SELECT page_id, name FROM page WHERE username IN (%S)", []string{"Starbucks", "McDonalds", "Dominos"})
	err := fql.Exec()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if fql.Result.Len() != 3 {
		t.Errorf("Bad length: %d", fql.Result.Len())
		t.Fail()
	}
}

func TestArrayCombinedQuery(t *testing.T) {
	fql := NewQuery("SELECT page_id, name, fan_count FROM page WHERE username IN (%S) AND fan_count > %d", []string{"coca-cola", "pepsi", "burgerking", "McDonalds"}, int64(3e7))
	err := fql.Exec()
	if err != nil {
		t.Errorf("Error executing fql: %s", err)
		t.Fail()
	}

	if fql.Result.Len() != 3 {
		t.Errorf("Unexpected number of results: %d", fql.Result.Len())
		t.Fail()
	}
}

func TestBadQuery(t *testing.T) {
	// invalid field in query
	fql := NewQuery("SELECT id, name, fan_count FROM page WHERE username IN ('pepsi') AND fan_count > %d", int64(3e7))
	err := fql.Exec()
	switch err.(type) {
	case Error:
	default:
		t.Errorf("Expected an error of type fql.Error, got: %#v", err)
		t.Fail()
	}
}
