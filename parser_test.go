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
	pkgName := "github.com/ConorNevin/traceable/internal/tests/geometry"
	pkgs, err := packages.Load(cfg, pkgName)
	c.Assert(err, qt.IsNil)
	c.Assert(pkgs, qt.HasLen, 1)

	pkg, err := pp.parsePackage(pkgs[0])
	c.Assert(err, qt.IsNil)
	c.Check(pkg.name, qt.Equals, "geometry")
	c.Assert(pkg.interfaces, qt.HasLen, 1)

	pm := getPackageMap(pp)

	i := pkg.interfaces[0]
	c.Check(i.name, qt.Equals, "Geometry")
	c.Assert(i.methods, qt.HasLen, 1)
	c.Check(i.methods[0].name, qt.Equals, "Area")
	c.Assert(i.methods[0].args, qt.HasLen, 1)
	c.Check(i.methods[0].args[0].String(pm, pkgName), qt.Equals, "context.Context")
	c.Assert(i.methods[0].returns, qt.HasLen, 2)
	c.Check(i.methods[0].returns[0].String(pm, pkgName), qt.Equals, "float64")
	c.Check(i.methods[0].returns[1].String(pm, pkgName), qt.Equals, "error")
}

func getPackageMap(p *parser) map[string]string {
	pm := make(map[string]string, len(p.imports))
	for _, ip := range p.imports {
		pm[ip.Path] = ip.Name
	}
	return pm
}
