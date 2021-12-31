package traceable

import (
	"go/ast"
	"testing"

	qt "github.com/frankban/quicktest"
	"golang.org/x/tools/go/packages"
)

func Test_parser_parsePackage(t *testing.T) {
	pp := &parser{
		imports:            make(map[string]ImportedPackage),
		importedInterfaces: make(map[string]map[string]*ast.InterfaceType),
		otherInterfaces:    make(map[string]map[string]*ast.InterfaceType),
	}

	c := qt.New(t)

	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedTypesInfo |
			packages.NeedSyntax |
			packages.NeedTypes,
	}
	pkgs, err := packages.Load(cfg, "github.com/ConorNevin/traceable/internal/tests/geometry")
	c.Assert(err, qt.IsNil)
	c.Assert(pkgs, qt.HasLen, 1)

	pkg, err := pp.parsePackage(pkgs[0])
	c.Assert(err, qt.IsNil)
	c.Check(pkg.name, qt.Equals, "geometry")
	c.Assert(pkg.interfaces, qt.HasLen, 1)

	i := pkg.interfaces[0]
	c.Check(i.name, qt.Equals, "Geometry")
	c.Assert(i.methods, qt.HasLen, 1)
	c.Check(i.methods[0].name, qt.Equals, "Area")
	c.Assert(i.methods[0].args, qt.HasLen, 1)
	c.Check(i.methods[0].args[0].String(), qt.Equals, "context.Context")
	c.Assert(i.methods[0].returns, qt.HasLen, 2)
	c.Check(i.methods[0].returns[0].String(), qt.Equals, "float64")
	c.Check(i.methods[0].returns[1].String(), qt.Equals, "error")
}
