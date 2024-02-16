package sqlbuild

import (
	"errors"
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

			if errors.Is(tC.wantError, err) {
				t.Errorf("Wanted error=%s, got %s", tC.wantError, err)
			}

			if tC.wantError == nil && tC.wantName != sName {
				t.Errorf("Wanted struct name '%s', got '%s'", tC.wantName, sName)
			}
		})
	}
}
