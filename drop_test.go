package sqlbuild

import "testing"

func TestDrop(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	p := Person{}

	want := "drop table Person"
	q, err := Drop(p)
	if err != nil {
		t.Error(err)
	}

	if q != want {
		t.Errorf("Got '%s', want '%s'", q, want)
	}
}
