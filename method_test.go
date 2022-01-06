package traceable

import (
	"go/token"
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
			name:   "no arguments",
			method: Method{},
			want:   false,
		},
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

func Test_Method_imports(t *testing.T) {
	tests := []struct {
		name   string
		method Method
		want   []string
	}{
		{
			name:   "no imports",
			method: Method{},
			want:   nil,
		},
		{
			name: "uses only built in types",
			method: Method{
				args: []types.Type{
					types.Typ[types.Bool],
				},
			},
			want: nil,
		},
		{
			name: "arg imports standard library package",
			method: Method{
				args: []types.Type{
					newContextType(),
				},
			},
			want: []string{"context"},
		},
		{
			name: "returns imports standard library package",
			method: Method{
				returns: []types.Type{
					newContextType(),
					newErrorType(),
				},
			},
			want: []string{"context"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := keys(tt.method.imports())
			if len(tt.want) == 0 {
				qt.Check(t, got, qt.HasLen, 0)
			} else {
				qt.Check(t, got, qt.ContentEquals, tt.want)
			}
		})
	}
}

func Test_importsOf(t *testing.T) {
	tests := []struct {
		name    string
		typ     types.Type
		imports []string
	}{
		{
			name:    "pointer",
			typ:     types.NewPointer(newContextType()),
			imports: []string{"context"},
		},
		{
			name:    "slice",
			typ:     types.NewSlice(newContextType()),
			imports: []string{"context"},
		},
		{
			name:    "array",
			typ:     types.NewArray(newContextType(), 5),
			imports: []string{"context"},
		},
		{
			name:    "chan",
			typ:     types.NewChan(types.SendRecv, newContextType()),
			imports: []string{"context"},
		},
		{
			name: "map",
			typ: types.NewMap(
				newContextType(),
				types.NewNamed(
					types.NewTypeName(
						token.NoPos,
						types.NewPackage("net/http", "http"),
						"Request",
						types.NewStruct(nil, nil),
					),
					types.NewStruct(nil, nil),
					nil,
				),
			),
			imports: []string{"context", "net/http"},
		},
		{
			name: "signature",
			typ: types.NewSignature(
				nil,
				types.NewTuple(
					types.NewParam(
						token.NoPos,
						types.NewPackage("context", "context"),
						"",
						newContextType(),
					),
				),
				types.NewTuple(
					types.NewParam(
						token.NoPos,
						types.NewPackage("net/http", "http"),
						"",
						newType("net/http", "Request"),
					),
				),
				false,
			),
			imports: []string{"context", "net/http"},
		},
		{
			name: "named type local to Interface",
			typ: types.NewNamed(
				types.NewTypeName(
					token.NoPos,
					nil,
					"Request",
					types.NewStruct(nil, nil),
				),
				types.NewStruct(nil, nil),
				nil,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := keys(importsOf(tt.typ))
			if len(tt.imports) == 0 {
				qt.Check(t, got, qt.HasLen, 0)
				return
			}

			qt.Check(t, got, qt.ContentEquals, tt.imports)
		})
	}
}

func newType(pkg, name string) types.Type {
	pkgs, _ := packages.Load(&packages.Config{
		Mode: packages.NeedTypes | packages.NeedImports,
	}, pkg)

	return pkgs[0].Types.Scope().Lookup(name).Type()
}

func keys(m map[string]struct{}) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
