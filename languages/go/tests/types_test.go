package oso

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	oso "github.com/osohq/oso/languages/go/pkg"
)

func TestSerialize(t *testing.T) {
	expected := `{"Number":{"Integer":123}}`
	int123 := oso.NumericInteger(123)
	term := oso.Value{&oso.ValueNumber{
		&int123,
	}}
	s, err := json.Marshal(term)
	if err != nil {
		t.Fatal(err)
	}
	if string(s) != expected {
		t.Fatal(fmt.Errorf("expected %#v got %#v", expected, string(s)))
	}

}

func TestDeserialize(t *testing.T) {
	jsonTerm := []byte(`{
        "Call": {
            "name": "foo",
            "args": [{"Number": {"Integer": 0}}],
            "kwargs": {"bar": {"Number": {"Integer": 1}}}
        }
	}`)

	var term oso.Value
	err := json.Unmarshal(jsonTerm, &term)
	if err != nil {
		t.Fatal(err)
	}
	int0 := oso.NumericInteger(0)
	int1 := oso.NumericInteger(1)
	expectedCall := oso.ValueCall{
		Name:   "foo",
		Args:   []oso.Value{{&oso.ValueNumber{&int0}}},
		Kwargs: &map[string]oso.Value{"bar": {&oso.ValueNumber{&int1}}},
	}
	expected := oso.Value{&expectedCall}
	if !cmp.Equal(term, expected) || !reflect.DeepEqual(term, expected) {
		t.Error(fmt.Errorf("Result differs from expected:\n%s", cmp.Diff(term, expected)))
	}

	jsonErrTerm := []byte(`{"kind":{"Parse":{"InvalidTokenCharacter":{"token":"this is not","c":"\n","loc":24}}},"formatted":"'\\n' is not a valid character. Found in this is not at line 1, column 25"}`)
	var errTerm oso.FormattedPolarError
	err = json.Unmarshal(jsonErrTerm, &errTerm)
	if err != nil {
		t.Fatal(err)
	}
	expectedErr := oso.FormattedPolarError{
		Kind: oso.ErrorKind{
			&oso.ErrorKindParse{
				&oso.ParseErrorInvalidTokenCharacter{
					Token: "this is not",
					C:     "\n",
					Loc:   24,
				},
			},
		},
		Formatted: "'\\n' is not a valid character. Found in this is not at line 1, column 25",
	}
	if !cmp.Equal(errTerm, expectedErr) {
		t.Error(fmt.Errorf("Result differs from expected:\n%s", cmp.Diff(errTerm, expectedErr)))
	}

}