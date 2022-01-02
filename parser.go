package traceable

import (
	"errors"
	"fmt"
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"
	"log"
	"path"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

type parser struct {
	imports            map[string]ImportedPackage
	importedInterfaces map[string]map[string]*ast.InterfaceType

	otherInterfaces map[string]map[string]*ast.InterfaceType
}

func (p *parser) parsePackage(pkg *packages.Package) (*Package, error) {
	var interfaces []*Interface
	for _, file := range pkg.Syntax {
		is, err := p.parseFile(pkg.PkgPath, file)
		if err != nil {
			return nil, err
		}

		interfaces = append(interfaces, is...)
	}

	return &Package{
		name:       pkg.Name,
		importPath: pkg.PkgPath,
		interfaces: interfaces,
	}, nil
}

func (p *parser) parseFile(path string, f *ast.File) ([]*Interface, error) {
	if err := p.parseImports(f); err != nil {
		return nil, err
	}

	p.otherInterfaces[path] = getInterfaces(f)

	var is []*Interface
	for name, typ := range p.otherInterfaces[path] {
		i, err := p.parseInterface(name, path, typ)
		if err != nil {
			return nil, err
		}
		is = append(is, i)

		i.imports = p.imports
	}

	return is, nil
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
	case *ast.ChanType:
		t, err := p.parseType(pkg, ft.Value)
		if err != nil {
			return nil, err
		}
		t.isChan = true
		switch ft.Dir {
		case ast.SEND:
			t.chanDir = ChanSend
		case ast.RECV:
			t.chanDir = ChanRecv
		}

		return t, nil
	case *ast.Ellipsis:
		t, err := p.parseType(pkg, ft.Elt)
		if err != nil {
			return nil, err
		}
		t.isVariadic = true
		return t, nil
	case *ast.FuncType:
		in, out, err := p.parseFunc(pkg, ft)
		if err != nil {
			return nil, err
		}

		return &Type{
			isFunc:        true,
			funcArguments: in,
			funcReturns:   out,
		}, nil
	case *ast.Ident:
		var t Type
		t.value = ft.Name

		if ft.IsExported() {
			maybeImportedPkg, ok := p.imports[pkg]
			if ok {
				pkg = maybeImportedPkg.Name
			}

			t.pkg = pkg
		}

		return &t, nil
	case *ast.InterfaceType:
		if ft.Methods != nil && len(ft.Methods.List) > 0 {
			return nil, fmt.Errorf("can not handle non-empty unnamed interfaces")
		}
		return &Type{value: "interface{}"}, nil
	case *ast.SelectorExpr:
		pkgName := ft.X.(*ast.Ident).String()
		pkg, ok := p.imports[pkgName]
		if !ok {
			return nil, fmt.Errorf("unknown package %q", pkgName)
		}
		return &Type{pkg: pkg.Path, value: ft.Sel.String()}, nil
	case *ast.MapType:
		key, err := p.parseType(pkg, ft.Key)
		if err != nil {
			return nil, err
		}
		value, err := p.parseType(pkg, ft.Value)
		if err != nil {
			return nil, err
		}

		return &Type{
			mapKey:   key,
			mapValue: value,
			isMap:    true,
		}, nil
	case *ast.StarExpr:
		t, err := p.parseType(pkg, ft.X)
		if err != nil {
			return nil, err
		}
		t.isPointer = true
		return t, nil
	case *ast.StructType:
		if ft.Fields != nil && len(ft.Fields.List) > 0 {
			return nil, errors.New("unable to handle non-empty unnamed structs")
		}
		return &Type{value: "struct{}"}, nil
	case *ast.ParenExpr:
		return p.parseType(pkg, ft.X)
	default:
		log.Fatalf("internal error: unexpected type: %T", typ)
		return nil, nil
	}
}

func (p *parser) parseFunc(pkg string, f *ast.FuncType) ([]*Parameter, []*Parameter, error) {
	var in, out []*Parameter
	var err error
	if f.Params != nil {
		in, err = p.parseFieldList(pkg, f.Params.List)
		if err != nil {
			return nil, nil, err
		}
	}
	if f.Results != nil {
		out, err = p.parseFieldList(pkg, f.Results.List)
		if err != nil {
			return nil, nil, err
		}
	}

	return in, out, nil
}

func (p *parser) parseInterface(name, pkg string, f *ast.InterfaceType) (*Interface, error) {
	i := Interface{name: name}
	for _, f := range f.Methods.List {
		switch v := f.Type.(type) {
		case *ast.FuncType:
			if n := len(f.Names); n != 1 {
				return nil, fmt.Errorf("expected one name for interface %q, got %d", i.name, n)
			}

			m := Method{
				name: f.Names[0].String(),
			}

			// if we already have a function with this name then we want
			// this one to supersede it.
			removeIdx := -1
			for idx, im := range i.methods {
				if im.name == m.name {
					removeIdx = idx
				}
			}
			if removeIdx != -1 {
				i.methods = append(i.methods[:removeIdx], i.methods[removeIdx+1:]...)
			}

			var err error
			m.args, m.returns, err = p.parseFunc(pkg, v)
			if err != nil {
				return nil, err
			}

			i.methods = append(i.methods, m)
		case *ast.Ident:
			embedType := p.otherInterfaces[pkg][v.String()]
			if embedType == nil {
				embedType = p.importedInterfaces[pkg][v.String()]
			}

			var embed *Interface
			if embedType != nil {
				var err error
				embed, err = p.parseInterface(v.String(), pkg, embedType)
				if err != nil {
					return nil, err
				}
			} else {
				if v.String() == "error" {
					embed = &Interface{
						name: "error",
						methods: []Method{
							{
								name: "Error",
								returns: []*Parameter{
									{
										typ: &Type{
											value: "string",
										},
									},
								},
							},
						},
					}
				} else {
					return nil, fmt.Errorf("unknown embedded interface %s", v.String())
				}
			}

			for _, m := range embed.methods {
				if i.hasMethod(m) {
					continue
				}

				i.methods = append(i.methods, m)
			}
		case *ast.SelectorExpr:
			// Embedded interface in another package.
			filePkg, sel := v.X.(*ast.Ident).String(), v.Sel.String()
			embeddedPkg, ok := p.imports[filePkg]
			if !ok {
				return nil, fmt.Errorf("unknown package %s", filePkg)
			}

			var embeddedIface *Interface
			var err error
			embeddedIfaceType := p.importedInterfaces[filePkg][sel]
			if embeddedIfaceType != nil {
				embeddedIface, err = p.parseInterface(sel, filePkg, embeddedIfaceType)
				if err != nil {
					return nil, err
				}
			} else {
				path := embeddedPkg.Path
				embeddedparser, err := p.createNewParser(path)
				if err != nil {
					return nil, err
				}

				if embeddedIfaceType = embeddedparser.importedInterfaces[path][sel]; embeddedIfaceType == nil {
					return nil, fmt.Errorf("unknown embedded interface %s.%s", path, sel)
				}
				embeddedIface, err = embeddedparser.parseInterface(sel, path, embeddedIfaceType)
				if err != nil {
					return nil, err
				}
			}

			for _, m := range embeddedIface.methods {
				if i.hasMethod(m) {
					continue
				}

				i.methods = append(i.methods, m)
			}
		default:
			return nil, fmt.Errorf("unable to handle method of type %T", f.Type)
		}
	}

	return &i, nil
}

func (p *parser) parseImports(file *ast.File) error {
	p.imports = make(map[string]ImportedPackage)
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
			p.imports[pkgName] = ImportedPackage{
				Name: pkgName,
				Path: importPath,
			}
		}
	}

	return nil
}

func (p *parser) createNewParser(path string) (*parser, error) {
	log.Printf("building parser for %s", path)
	n := &parser{
		imports:            make(map[string]ImportedPackage),
		importedInterfaces: make(map[string]map[string]*ast.InterfaceType),
		otherInterfaces:    make(map[string]map[string]*ast.InterfaceType),
	}

	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedTypesInfo |
			packages.NeedSyntax |
			packages.NeedTypes,
	}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			if _, ok := n.importedInterfaces[path]; !ok {
				n.importedInterfaces[path] = make(map[string]*ast.InterfaceType)
			}
			for name, it := range getInterfaces(file) {
				n.importedInterfaces[path][name] = it
			}

			if err := n.parseImports(file); err != nil {
				return nil, err
			}
		}
	}

	return n, nil
}

func getInterfaces(f *ast.File) map[string]*ast.InterfaceType {
	interfaces := make(map[string]*ast.InterfaceType)
	ast.Inspect(f, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		it, ok := ts.Type.(*ast.InterfaceType)
		if !ok {
			return false
		}

		log.Printf("found %s in %s", ts.Name.Name, f.Name.Name)
		interfaces[ts.Name.Name] = it
		return true
	})

	return interfaces
}
