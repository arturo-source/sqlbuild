package sqlbuild

import (
	"errors"
	"testing"
)

func TestCreate(t *testing.T) {
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

	want := `create table "Person" ("Id" int not null auto_increment primary key, "Name" text not null, "Age" int)`
	q, err := Create(p)
	if err != nil {
		t.Error(err)
	}

	if q != want {
		t.Errorf("Want query '%s', got query '%s'", want, q)
	}
}

func TestCreateUsingInvalidType(t *testing.T) {
	type Person struct {
		Name     string
		Age      int
		Measures struct {
			Height float64
			Weight float64
		}
	}
	p := Person{}

	wantErr := ErrNoValidType{"struct"}
	_, err := Create(p)
	if !errors.Is(err, wantErr) {
		t.Errorf("Got '%s', want '%s'", err, wantErr)
	}
}
