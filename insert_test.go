package sqlbuild

import (
	"reflect"
	"strings"
	"testing"
)

func TestInsert(t *testing.T) {
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

	q, err := Insert(p)
	if err != nil {
		t.Error(err)
	}

	wantPrefix := "insert into 'Person' "
	if !strings.HasPrefix(q, wantPrefix) {
		t.Errorf("Want prefix '%s', got query '%s'", wantPrefix, q)
	}

	wantContaining := ") values ("
	if !strings.Contains(q, wantContaining) {
		t.Errorf("Want containing '%s', got query '%s'", wantContaining, q)
	}

	commas := strings.Count(q, ",")
	wantCommas := 4
	if commas != wantCommas {
		t.Errorf("Got %d commas, want %d '%s'", commas, wantCommas, q)
	}

	fields := newFields(reflect.ValueOf(p))
	for f := range fields.nameValues {
		if !strings.Contains(q, f) {
			t.Errorf("Want containing %s, got query '%s'", f, q)
		}
	}
}

func TestInsertMultiple(t *testing.T) {
	// TODO
}
