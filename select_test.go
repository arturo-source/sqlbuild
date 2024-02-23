package sqlbuild

import (
	"testing"
)

func TestSelectAll(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	p := Person{}

	want := `SELECT * FROM "Person"`
	q, err := SelectAll(p)
	if err != nil {
		t.Error(err)
	}

	if q != want {
		t.Errorf("Got '%s', want '%s'", q, want)
	}
}

func TestSelectById(t *testing.T) {
	type Person struct {
		Id   int `db:"id"`
		Name string
		Age  int
	}
	p := Person{Id: 10}

	want := `SELECT * FROM "Person" WHERE "id" = 10`
	q, err := SelectById(p)
	if err != nil {
		t.Error(err)
	}

	if q != want {
		t.Errorf("Got '%s', want '%s'", q, want)
	}
}

func TestSelectByIdString(t *testing.T) {
	type Person struct {
		Id   string `db:"id"`
		Name string
		Age  int
	}
	p := Person{Id: "10"}

	want := `SELECT * FROM "Person" WHERE "id" = '10'`
	q, err := SelectById(p)
	if err != nil {
		t.Error(err)
	}

	if q != want {
		t.Errorf("Got '%s', want '%s'", q, want)
	}
}
