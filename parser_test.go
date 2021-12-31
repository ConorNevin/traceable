package traceable

import (
	"go/ast"
	goparser "go/parser"
	"go/token"
	"testing"

	qt "github.com/frankban/quicktest"
)

func Test_parser_parsePackage(t *testing.T) {
	pp := &parser{
		imports:            make(map[string]ImportedPackage),
		importedInterfaces: make(map[string]map[string]*ast.InterfaceType),
		otherInterfaces:    make(map[string]map[string]*ast.InterfaceType),
	}

	fs := token.NewFileSet()
	file, err := goparser.ParseFile(fs, "internal/testdata/geometry/geometry.go", nil, 0)
	qt.Assert(t, err, qt.IsNil)

	interfaces, err := pp.parseFile("", file)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, interfaces, qt.HasLen, 1)

	i := interfaces[0]
	qt.Check(t, i.name, qt.Equals, "Geometry")
	qt.Assert(t, i.methods, qt.HasLen, 1)
	qt.Check(t, i.methods[0].name, qt.Equals, "Area")
	qt.Assert(t, i.methods[0].args, qt.HasLen, 1)
	qt.Check(t, i.methods[0].args[0].String(), qt.Equals, "context.Context")
	qt.Assert(t, i.methods[0].returns, qt.HasLen, 2)
	qt.Check(t, i.methods[0].returns[0].String(), qt.Equals, "float64")
	qt.Check(t, i.methods[0].returns[1].String(), qt.Equals, "error")
}
