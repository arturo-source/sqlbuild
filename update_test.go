package sqlbuild

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestUpdate(t *testing.T) {
	type Person struct {
		Id   int
		Name string
		Age  int
	}
	p := Person{
		Id:   1,
		Name: "John",
		Age:  10,
	}

	q, err := Update(p)
	if err != nil {
		t.Error(err)
	}

	wantPrefix := "update 'Person' set "
	if !strings.HasPrefix(q, wantPrefix) {
		t.Errorf("Want prefix '%s', got query '%s'", wantPrefix, q)
	}

	wantSuffix := " where 'Id' = 1"
	if !strings.HasSuffix(q, wantSuffix) {
		t.Errorf("Want suffix '%s', got query '%s'", wantSuffix, q)
	}

	commas := strings.Count(q, ",")
	wantCommas := 2
	if commas != wantCommas {
		t.Errorf("Got %d commas, want %d '%s'", commas, wantCommas, q)
	}

	fields := getStructFields(reflect.ValueOf(p))
	for f := range fields {
		if !strings.Contains(q, f) {
			t.Errorf("Want containing %s, got query '%s'", f, q)
		}
	}
}

func TestUpdateWithoutId(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	p := Person{}

	_, err := Update(p)
	if !errors.Is(err, ErrNoId) {
		t.Errorf("Got '%s', want '%s'", err, ErrNoId)
	}
}
