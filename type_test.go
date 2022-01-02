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
		{
			name: "pointer type",
			typ: Type{
				isPointer: true,
				value:     "float64",
			},
			want: "*float64",
		},
		{
			name: "slice of pointers",
			typ: Type{
				isSlice:   true,
				isPointer: true,
				value:     "string",
			},
			want: "[]*string",
		},
		{
			name: "array",
			typ: Type{
				isSlice:     true,
				value:       "string",
				arrayLength: 2,
			},
			want: "[2]string",
		},
		{
			name: "map of basic types",
			typ: Type{
				isMap: true,
				mapKey: &Type{
					value: "string",
				},
				mapValue: &Type{
					isSlice: true,
					pkg:     "context",
					value:   "Context",
				},
			},
			want: "map[string][]context.Context",
		},
		{
			name: "channel",
			typ: Type{
				isChan: true,
				value:  "string",
			},
			want: "chan string",
		},
		{
			name: "channel receiver",
			typ: Type{
				isChan:  true,
				chanDir: ChanRecv,
				value:   "string",
			},
			want: "<-chan string",
		},
		{
			name: "channel sender",
			typ: Type{
				isChan:  true,
				chanDir: ChanSend,
				value:   "string",
			},
			want: "chan<- string",
		},
		{
			name: "function with arguments and return values",
			typ: Type{
				isFunc: true,
				funcArguments: []*Parameter{
					{
						typ: &Type{
							pkg:   "context",
							value: "Context",
						},
					},
					{
						typ: &Type{
							value: "string",
						},
					},
				},
				funcReturns: []*Parameter{
					{
						typ: &Type{
							value: "error",
						},
					},
				},
			},
			want: "func(context.Context, string) error",
		},
		{
			name: "function with no return values",
			typ: Type{
				isFunc: true,
				funcArguments: []*Parameter{
					{
						typ: &Type{
							pkg:   "context",
							value: "Context",
						},
					},
					{
						typ: &Type{
							value: "string",
						},
					},
				},
			},
			want: "func(context.Context, string)",
		},
		{
			name: "function with no arguments",
			typ: Type{
				isFunc: true,
				funcReturns: []*Parameter{
					{
						typ: &Type{
							pkg:   "context",
							value: "Context",
						},
					},
					{
						typ: &Type{
							value: "string",
						},
					},
				},
			},
			want: "func() (context.Context, string)",
		},
		{
			name: "function with no arguments or return values",
			typ: Type{
				isFunc: true,
			},
			want: "func()",
		},
		{
			name: "variadic options",
			typ: Type{
				isVariadic: true,
				value:      "string",
			},
			want: "...string",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t *testing.T) {
			qt.Check(t, tt.typ.String(pm, override), qt.Equals, tt.want)
		})
	}
}
