package traceable

import (
	"go/ast"
	"go/token"
	"log"
)

type Interface struct {
	name       string
	importPath string
	methods    []Method
	imports    map[string]ImportedPackage
}

func (i Interface) hasMethod(m Method) bool {
	for _, im := range i.methods {
		if im.name == m.name {
			return true
		}
	}

	return false
}

func (g *Generator) findInterface(n ast.Node) bool {
	inter, ok := n.(*ast.GenDecl)
	if !ok || inter.Tok != token.TYPE {
		// We only care about interfaces.
		return true
	}

	name := g.Interface.name

	typ := ""
	for _, spec := range inter.Specs {
		tspec := spec.(*ast.TypeSpec)
		if tspec.Name != nil {
			typ = tspec.Name.Name
		}

		// Check if this is the interface we're looking for and
		// skip if not.
		if typ != name {
			log.Printf("skipping %s (looking for %s)", typ, name)
			continue
		}

		itype, ok := tspec.Type.(*ast.InterfaceType)
		if !ok {
			log.Fatalf("%s is an %T, not an interface,", typ, tspec.Type)
		}

		if err := g.setMethods(itype); err != nil {
			log.Fatal(err)
		}

		return false
	}

	return false
}

func (g *Generator) setMethods(inter *ast.InterfaceType) error {
	log.Printf("setting methods (%d)", len(inter.Methods.List))
	for _, fl := range inter.Methods.List {
		switch typ := fl.Type.(type) {
		case *ast.FuncType:
			args, err := g.parser.parseFieldList(g.Interface.importPath, typ.Params.List)
			if err != nil {
				return err
			}

			results, err := g.parser.parseFieldList(g.Interface.importPath, typ.Results.List)
			if err != nil {
				return err
			}

			g.Interface.methods = append(g.Interface.methods, Method{
				name:    fl.Names[0].Name,
				args:    args,
				returns: results,
			})
		case *ast.Ident:
		case *ast.SelectorExpr:
		default:
			log.Fatalf("unexpected type: %T", fl.Type)
		}
	}

	return nil
}
