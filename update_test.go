package sqlbuild

import (
	"errors"
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

	want := `UPDATE "Person" SET "Id" = 1, "Name" = 'John', "Age" = 10 WHERE "Id" = 1`
	q, err := Update(p)
	if err != nil {
		t.Error(err)
	}

	if q != want {
		t.Errorf("Want query '%s', got query '%s'", want, q)
	}
}

func TestUpdateWithFieldPointer(t *testing.T) {
	type Person struct {
		Id   int
		Name string
		Age  *int
	}
	p := Person{
		Id:   1,
		Name: "John",
		Age:  nil,
	}

	want := `UPDATE "Person" SET "Id" = 1, "Name" = 'John', "Age" = null WHERE "Id" = 1`
	q, err := Update(p)
	if err != nil {
		t.Error(err)
	}

	if q != want {
		t.Errorf("Want query '%s', got query '%s'", want, q)
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
