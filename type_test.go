package traceable

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestType_String(t1 *testing.T) {
	pm := map[string]string{
		"context":             "context",
		"this/is/nested":      "nested",
		"this/is/also/nested": "nested1",
		"the/root":            "root",
	}
	override := "the/root"

	tests := []struct {
		name string
		typ  Type
		want string
	}{
		{
			name: "basic type",
			typ: Type{
				value: "float64",
			},
			want: "float64",
		},
		{
			name: "with imported package",
			typ: Type{
				pkg:   "this/is/nested",
				value: "foo",
			},
			want: "nested.foo",
		},
		{
			name: "in override package",
			typ: Type{
				pkg:   override,
				value: "foo",
			},
			want: "foo",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t *testing.T) {
			qt.Check(t, tt.typ.String(pm, override), qt.Equals, tt.want)
		})
	}
}
