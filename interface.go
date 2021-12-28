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

	parser *parser
}

func (f *Interface) findInterfaces(n ast.Node) bool {
	inter, ok := n.(*ast.GenDecl)
	if !ok || inter.Tok != token.TYPE {
		// We only care about interfaces.
		return true
	}

	typ := ""
	for _, spec := range inter.Specs {
		tspec := spec.(*ast.TypeSpec)
		if tspec.Name != nil {
			typ = tspec.Name.Name
		}

		// Check if this is the interface we're looking for and
		// skip if not.
		if typ != f.name {
			log.Printf("skipping %s (looking for %s)", typ, f.name)
			continue
		}

		itype, ok := tspec.Type.(*ast.InterfaceType)
		if !ok {
			log.Fatalf("%s is an %T, not an interface,", typ, tspec.Type)
		}

		if err := f.setMethods(itype); err != nil {
			log.Fatal(err)
		}
	}
	return false
}

func (f *Interface) setMethods(inter *ast.InterfaceType) error {
	log.Printf("setting methods (%d)", len(inter.Methods.List))
	for _, fl := range inter.Methods.List {
		typ, ok := fl.Type.(*ast.FuncType)
		if !ok {
			log.Printf("unexpected type: %T", typ)
		}

		args, err := f.parser.parseFieldList(f.importPath, typ.Params.List)
		if err != nil {
			return err
		}

		results, err := f.parser.parseFieldList(f.importPath, typ.Results.List)
		if err != nil {
			return err
		}

		f.methods = append(f.methods, Method{
			name:    fl.Names[0].Name,
			args:    args,
			returns: results,
		})
	}

	return nil
}
