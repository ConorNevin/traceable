package traceable

import (
	"go/ast"
	"go/types"
	"log"

	"golang.org/x/tools/go/packages"
)

type parser struct {
	imports            map[string]ImportedPackage
	importedInterfaces map[string]map[string]*ast.InterfaceType

	otherInterfaces map[string]map[string]*ast.InterfaceType
}

func (p *parser) parsePackage(pkg *packages.Package) (*Package, error) {
	var interfaces []*Interface

	scope := pkg.Types.Scope()
	for _, name := range scope.Names() {
		o := scope.Lookup(name)

		switch ti := o.Type().Underlying().(type) {
		case *types.Interface:
			log.Printf("found interface: %s\n", o.Name())
			i, err := p.parseInterface(o.Name(), pkg.PkgPath, ti)
			if err != nil {
				return nil, err
			}

			interfaces = append(interfaces, i)
		}
	}

	return &Package{
		name:       pkg.Name,
		importPath: pkg.PkgPath,
		imports:    pkg.Types.Imports(),
		interfaces: interfaces,
	}, nil
}

func (p *parser) parseInterface(name, pkg string, ti *types.Interface) (*Interface, error) {
	i := Interface{name: name, methods: make([]Method, ti.NumMethods())}
	for idx := 0; idx < ti.NumMethods(); idx++ {
		m, err := p.parseFunc(ti.Method(idx))
		if err != nil {
			return nil, err
		}

		i.methods[idx] = *m
	}

	for idx := 0; idx < ti.NumEmbeddeds(); idx++ {
		e := ti.EmbeddedType(idx)
		log.Printf("found embedded type with type - %T", e)
	}

	return &i, nil
}

func (p *parser) parseFunc(f *types.Func) (*Method, error) {
	sig := f.Type().(*types.Signature)
	m := &Method{
		name:       f.Name(),
		args:       make([]types.Type, sig.Params().Len()),
		returns:    make([]types.Type, sig.Results().Len()),
		isVariadic: sig.Variadic(),
	}

	for i := range m.args {
		m.args[i] = sig.Params().At(i).Type()
	}
	for i := range m.returns {
		m.returns[i] = sig.Results().At(i).Type()
	}

	return m, nil
}
