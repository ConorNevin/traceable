package traceable

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestInterface_hasMethod(t *testing.T) {
	tests := []struct {
		name    string
		m       Method
		methods []Method
		want    bool
	}{
		{
			name: "does not have method",
			m:    Method{name: "Foo"},
			methods: []Method{
				{
					name: "Bar",
				},
			},
		},
		{
			name: "has method",
			m:    Method{name: "Foo"},
			methods: []Method{
				{
					name: "Foo",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Interface{
				methods: tt.methods,
			}

			qt.Check(t, i.hasMethod(tt.m), qt.Equals, tt.want)
		})
	}
}
