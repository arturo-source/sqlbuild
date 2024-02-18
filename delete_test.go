package sqlbuild

import (
	"testing"
)

func TestDeleteAll(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	p := Person{}

	want := "delete from Person"
	q, err := DeleteAll(p)
	if err != nil {
		t.Error(err)
	}

	if q != want {
		t.Errorf("Got '%s', want '%s'", q, want)
	}
}

func TestDeleteById(t *testing.T) {
	type Person struct {
		Id   int `db:"id"`
		Name string
		Age  int
	}
	p := Person{Id: 10}

	want := "delete from 'Person' where 'id' = 10"
	q, err := DeleteById(p)
	if err != nil {
		t.Error(err)
	}

	if q != want {
		t.Errorf("Got '%s', want '%s'", q, want)
	}
}
