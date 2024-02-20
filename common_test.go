package sqlbuild

import (
	"errors"
	"reflect"
	"testing"
)

func TestGetStructFromPointer(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	p := Person{}
	val := 2

	testCases := []struct {
		desc      string
		gotStruct any
		wantError error
	}{
		{
			desc:      "A struct",
			gotStruct: p,
			wantError: nil,
		},
		{
			desc:      "Not a struct",
			gotStruct: val,
			wantError: ErrNoStruct,
		},
		{
			desc:      "A struct from a pointer",
			gotStruct: &p,
			wantError: nil,
		},
		{
			desc:      "Not a struct from a pointer",
			gotStruct: &val,
			wantError: ErrNoStruct,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := getStructFromPointer(tC.gotStruct)
			if !errors.Is(tC.wantError, err) {
				t.Errorf("Wanted error=%s, got %s", tC.wantError, err)
			}
		})
	}
}

func TestGetStructName(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	p := Person{}
	person, err := getStructFromPointer(p)
	if err != nil {
		t.Error(err)
	}

	want := "Person"
	pName := getStructName(person)
	if pName != want {
		t.Errorf("Got '%s', want '%s'", pName, want)
	}
}

func TestGetStructFieldNames(t *testing.T) {
	type PersonEmpty struct{}
	type Person struct {
		Name string
		Age  int
	}
	type Person2 struct {
		Name string `db:"name" json:"name-json"`
		Age  int    `db:"age"`
	}
	type Person3 struct {
		Name string `db:"name" json:"name-json"`
		Age  int
	}

	testCases := []struct {
		desc   string
		person any
		want   []string
	}{
		{
			desc:   "Empty struct",
			person: PersonEmpty{},
			want:   []string{},
		},
		{
			desc:   "Struct without tags",
			person: Person{},
			want:   []string{"Name", "Age"},
		},
		{
			desc:   "Struct with tags",
			person: Person2{},
			want:   []string{"name", "age"},
		},
		{
			desc:   "Struct with one tag",
			person: Person3{},
			want:   []string{"name", "Age"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			person, err := getStructFromPointer(tC.person)
			if err != nil {
				t.Error(err)
			}

			fields := newFields(person)
			if !reflect.DeepEqual(tC.want, fields.namesOrdered) {
				t.Errorf("Wanted %v, got %v", tC.want, fields.namesOrdered)
			}
		})
	}
}

func TestGetIdFromFields(t *testing.T) {
	emptyStruct := reflect.ValueOf(struct{}{})
	structWithId := reflect.ValueOf(struct{ Id int }{Id: 10})

	testCases := []struct {
		desc      string
		fields    reflect.Value
		wantKey   string
		wantValue any
		wantErr   error
	}{
		{
			desc:    "Error no id found",
			fields:  emptyStruct,
			wantErr: ErrNoId,
		},
		{
			desc:      "Id found",
			fields:    structWithId,
			wantKey:   "Id",
			wantValue: 10,
			wantErr:   nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			fields := newFields(tC.fields)
			key, value, err := fields.getId()
			if !errors.Is(tC.wantErr, err) {
				t.Errorf("Wanted '%s', got '%s'", tC.wantErr, err)
			}

			if tC.wantErr == nil && (tC.wantKey != key || tC.wantValue != value) {
				t.Errorf("Wanted %s = %v, got %s = %v", tC.wantKey, tC.wantValue, key, value)
			}
		})
	}
}

func TestSanitize(t *testing.T) {
	testCases := []struct {
		desc      string
		thing     any
		quoteMark string
		want      any
	}{
		{
			desc:      "Sanitize without quotation marks",
			thing:     "abc",
			quoteMark: "'",
			want:      "'abc'",
		},
		{
			desc:      "Sanitize with quotation marks",
			thing:     "a'bc",
			quoteMark: "'",
			want:      "'a''bc'",
		},
		{
			desc:      "Sanitize with quotation marks",
			thing:     10,
			quoteMark: "'",
			want:      10,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			sanitized := sanitize(tC.thing, tC.quoteMark)
			if tC.want != sanitized {
				t.Errorf("Wanted '%s', got '%s'", tC.want, sanitized)
			}
		})
	}
}
