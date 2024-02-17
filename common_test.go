package sqlbuild

import (
	"errors"
	"reflect"
	"testing"
)

func TestGetStructName(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	p := Person{}
	val := 2

	testCases := []struct {
		desc      string
		gotStruct any
		wantName  string
		wantError error
	}{
		{
			desc:      "A struct",
			gotStruct: p,
			wantName:  "Person",
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
			wantName:  "Person",
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
			sName, err := getStructName(tC.gotStruct)

			if !errors.Is(tC.wantError, err) {
				t.Errorf("Wanted error=%s, got %s", tC.wantError, err)
			}

			if tC.wantError == nil && tC.wantName != sName {
				t.Errorf("Wanted struct name '%s', got '%s'", tC.wantName, sName)
			}
		})
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

	isSameArray := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}

		var timesRepeated int
		for _, aa := range a {
			for _, bb := range b {
				if aa == bb {
					timesRepeated++
				}
			}
		}

		return timesRepeated == len(a)
	}
	getKeys := func(m Fields) []string {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}

		return keys
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
			fields := getStructFields(tC.person)
			fieldNames := getKeys(fields)

			if !isSameArray(tC.want, fieldNames) {
				t.Errorf("Wanted %v, got %v", tC.want, fieldNames)
			}
		})
	}
}

func TestGetIdFromFields(t *testing.T) {
	testCases := []struct {
		desc      string
		fields    Fields
		wantKey   string
		wantValue reflect.Value
		wantErr   error
	}{
		{
			desc:    "Error no id found",
			fields:  Fields{},
			wantErr: ErrNoId,
		},
		{
			desc:      "Id found",
			fields:    Fields{"Id": reflect.ValueOf(10)},
			wantKey:   "Id",
			wantValue: reflect.ValueOf(10),
			wantErr:   nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			key, value, err := getIdFromFields(tC.fields)
			if tC.wantErr != err {
				t.Errorf("Wanted '%s', got '%s'", tC.wantErr, err)
			}

			if tC.wantErr == nil && (tC.wantKey != key || tC.wantValue != value) {
				t.Errorf("Wanted %s = %v, got %s = %v", tC.wantKey, tC.wantValue, key, value)
			}
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	testCases := []struct {
		desc  string
		query string
		args  []any
		want  string
	}{
		{
			desc:  "Sanitize without quotation marks",
			query: "%s %v",
			args:  []any{"abc", reflect.ValueOf(10)},
			want:  "'abc' 10",
		},
		{
			desc:  "Sanitize with quotation marks",
			query: "%s %v",
			args:  []any{"a'bc", reflect.ValueOf(10)},
			want:  "'a''bc' 10",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			queryResult := sanitizeInput(tC.query, tC.args...)
			if tC.want != queryResult {
				t.Errorf("Wanted '%s', got '%s'", tC.want, queryResult)
			}
		})
	}
}
