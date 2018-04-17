package util

import (
	"reflect"
	"testing"
)

func TestUrandom(t *testing.T) {
	input := 123
	output := Urandom(input)

	typ := reflect.TypeOf(output)
	if typ.Kind() != reflect.Slice {
		t.Errorf("Not a slice: %s\n", typ.Kind())
	}

}
