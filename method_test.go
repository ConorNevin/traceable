package traceable

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func Test_Method_acceptsContext(t *testing.T) {
	tests := []struct {
		name   string
		method Method
		want   bool
	}{
		{
			name: "does not contain context",
			method: Method{
				args: []*Parameter{
					{
						typ: &Type{
							value: "string",
						},
					},
				},
			},
			want: false,
		},
		{
			name: "takes context as only argument",
			method: Method{
				args: []*Parameter{
					{
						typ: &Type{
							pkg:   "context",
							value: "Context",
						},
					},
				},
			},
			want: true,
		},
		{
			name: "takes context as one argument",
			method: Method{
				args: []*Parameter{
					{
						typ: &Type{
							pkg:   "context",
							value: "Context",
						},
					},
					{
						typ: &Type{
							pkg:   "http",
							value: "Request",
						},
					},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.method.acceptsContext()
			qt.Check(t, got, qt.Equals, tt.want)
		})
	}
}

func Test_Method_contextArg(t *testing.T) {
	tests := []struct {
		name   string
		method Method
		want   string
	}{
		{
			name:   "does not accept context as an argument",
			method: Method{},
			want:   "",
		},
		{
			name: "takes context as only argument",
			method: Method{
				args: []*Parameter{
					{
						typ: &Type{
							pkg:   "context",
							value: "Context",
						},
					},
				},
			},
			want: "a0",
		},
		{
			name: "takes context as one argument",
			method: Method{
				args: []*Parameter{
					{
						typ: &Type{
							pkg:   "context",
							value: "Context",
						},
					},
					{
						typ: &Type{
							pkg:   "http",
							value: "Request",
						},
					},
				},
			},
			want: "a0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.method.contextArg()
			qt.Check(t, got, qt.Equals, tt.want)
		})
	}
}
