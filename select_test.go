package sqlbuild

import "testing"

func TestSelectAll(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	p := Person{}

	want := "select * from Person"
	q, err := SelectAll(p)
	if err != nil {
		t.Error(err)
	}

	if q != want {
		t.Errorf("Got '%s', want '%s'", q, want)
	}
}
