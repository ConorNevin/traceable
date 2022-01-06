package traceable

import (
	"go/ast"
	"testing"

	qt "github.com/frankban/quicktest"
	"golang.org/x/tools/go/packages"
)

func Benchmark_parseFile(b *testing.B) {
	c := qt.New(b)

	pkgName := "internal/tests/performance/large"
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedTypesInfo |
			packages.NeedSyntax |
			packages.NeedTypes,
	}
	pkgs, err := packages.Load(cfg, pkgName)
	c.Assert(err, qt.IsNil)
	c.Assert(pkgs, qt.HasLen, 1)

	pp := &parser{
		imports:            make(map[string]ImportedPackage),
		importedInterfaces: make(map[string]map[string]*ast.InterfaceType),
		otherInterfaces:    make(map[string]map[string]*ast.InterfaceType),
	}

	for n := 0; n < b.N; n++ {
		_, _ = pp.parsePackage(pkgs[0])
	}
}

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
	t.Log("parsed")
	c.Assert(err, qt.IsNil)
	c.Check(pkg.name, qt.Equals, "geometry")
	c.Assert(pkg.interfaces, qt.HasLen, 1)

	t.Log("here")
	i := pkg.interfaces[0]
	c.Check(i.name, qt.Equals, "Geometry")
	c.Assert(i.methods, qt.HasLen, 2)
	t.Log("here again")
	c.Check(i.methods[0].name, qt.Equals, "Area")
	c.Assert(i.methods[0].args, qt.HasLen, 1)
	t.Log("starting...")
	c.Check(i.methods[0].args[0].String(), qt.Equals, "context.Context")
	t.Log("here")
	c.Assert(i.methods[0].returns, qt.HasLen, 2)
	c.Check(i.methods[0].returns[0].String(), qt.Equals, "float64")
	c.Check(i.methods[0].returns[1].String(), qt.Equals, "error")
}
