package traceable

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"
	"log"
	"path"
	"strconv"
	"strings"
)

type parser struct {
	imports map[string]ImportedPackage
}

func (p *parser) parseFieldList(pkg string, fields []*ast.Field) ([]*Parameter, error) {
	nf := 0
	for _, f := range fields {
		nn := len(f.Names)
		if nn == 0 {
			nn = 1
		}
		nf += nn
	}
	ps := make([]*Parameter, nf)
	i := 0

	for _, f := range fields {
		t, err := p.parseType(pkg, f.Type)
		if err != nil {
			return nil, err
		}

		if len(f.Names) == 0 {
			ps[i] = &Parameter{typ: t}
			i++
			continue
		}
		for _, name := range f.Names {
			ps[i] = &Parameter{name: name.Name, typ: t}
			i++
		}
	}

	return ps, nil
}

func (p *parser) parseType(pkg string, typ ast.Expr) (*Type, error) {
	switch ft := typ.(type) {
	case *ast.ArrayType:
		ln := -1
		if ft.Len != nil {
			var value string
			switch val := ft.Len.(type) {
			case *ast.BasicLit:
				value = val.Value
			case *ast.Ident:
				// when the length is a const defined locally
				value = val.Obj.Decl.(*ast.ValueSpec).Values[0].(*ast.BasicLit).Value
			case *ast.SelectorExpr:
				// when the length is a const defined in an external package
				usedPkg, err := importer.Default().Import(fmt.Sprintf("%s", val.X))
				if err != nil {
					log.Fatalf("unknown package in array length: %v", err)
				}
				ev, err := types.Eval(token.NewFileSet(), usedPkg, token.NoPos, val.Sel.Name)
				if err != nil {
					log.Fatalf("unknown constant in array length: %v", err)
				}
				value = ev.Value.String()
			}
			x, err := strconv.Atoi(value)
			if err != nil {
				log.Fatalf("bad array size: %v", err)
			}
			ln = x
		}
		t, err := p.parseType(pkg, ft.Elt)
		if err != nil {
			return t, err
		}
		t.isSlice = true
		t.arrayLength = ln

		return t, nil
	case *ast.Ident:
		var t Type
		t.value = ft.Name

		return &t, nil
	case *ast.SelectorExpr:
		pkgName := ft.X.(*ast.Ident).String()
		pkg, ok := p.imports[pkgName]
		if !ok {
			return nil, fmt.Errorf("unknown package %q", pkgName)
		}
		return &Type{packageName: pkg.Path, value: ft.Sel.String()}, nil
	default:
		log.Fatalf("internal error: unexpected type: %T", typ)
		return nil, nil
	}
}

func (p *parser) parseImports(file *ast.File) error {
	var importPaths []string
	for _, is := range file.Imports {
		if is.Name != nil {
			continue
		}
		importPath := is.Path.Value[1 : len(is.Path.Value)-1] // remove quotes
		importPaths = append(importPaths, importPath)
	}

	packagesName, err := createPackageMap(importPaths)
	if err != nil {
		return err
	}

	for _, is := range file.Imports {
		var pkgName string
		importPath := is.Path.Value[1 : len(is.Path.Value)-1]

		if is.Name != nil {
			if is.Name.Name == "_" {
				continue
			}
			pkgName = is.Name.Name
		} else {
			pkg, ok := packagesName[importPath]
			if !ok {
				// Fallback to import path suffix. Note that this is uncertain.
				_, last := path.Split(importPath)
				pkgName = strings.SplitN(last, ".", 2)[0]
			} else {
				pkgName = pkg
			}
		}

		if pkg, ok := p.imports[pkgName]; ok {
			p.imports[pkgName] = ImportedPackage{
				Name:       pkgName,
				Path:       pkg.Path,
				Duplicates: append(pkg.Duplicates, importPath),
			}
		} else {
			p.imports[pkgName] = ImportedPackage{Path: importPath}
		}
	}

	return nil
}
