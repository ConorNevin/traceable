package traceable

import (
	"go/types"
	"testing"

	qt "github.com/frankban/quicktest"
	"golang.org/x/tools/go/packages"
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
				args: []types.Type{
					newType("net/http", "Request"),
				},
			},
			want: false,
		},
		{
			name: "takes context as only argument",
			method: Method{
				args: []types.Type{
					newContextType(),
				},
			},
			want: true,
		},
		{
			name: "takes context as one argument",
			method: Method{
				args: []types.Type{
					newContextType(),
					newType("net/http", "Request"),
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
				args: []types.Type{
					newContextType(),
				},
			},
			want: "a0",
		},
		{
			name: "takes context as one argument",
			method: Method{
				args: []types.Type{
					newContextType(),
					newType("net/http", "Request"),
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

func newType(pkg, name string) types.Type {
	pkgs, _ := packages.Load(&packages.Config{
		Mode: packages.NeedTypes | packages.NeedImports,
	}, pkg)

	return pkgs[0].Types.Scope().Lookup(name).Type()
}
