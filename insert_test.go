package sqlbuild

import (
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

	want := "insert into 'Person' ('Id', 'Name', 'Age') values (1, 'John', 10)"
	q, err := Insert(p)
	if err != nil {
		t.Error(err)
	}

	if q != want {
		t.Errorf("Want query '%s', got query '%s'", want, q)
	}
}

func TestInsertMultiple(t *testing.T) {
	// TODO
}
