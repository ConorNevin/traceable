package traceable

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"
	"log"
	"strconv"
	"strings"
)

type Method struct {
	name    string
	args    []Type
	returns []Type
}

type Type struct {
	name       string
	value      string
	isPointer  bool
	isSlice    bool
	isVariadic bool

	arrayLength int
}

func newType(f *ast.Field) Type {
	t := parseType(f.Type)
	if len(f.Names) > 0 {
		t.name = f.Names[0].Name
	}

	return t
}

func parseType(typ ast.Expr) Type {
	switch ft := typ.(type) {
	case *ast.Ident:
		var t Type
		t.value = ft.Name

		return t
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
		t := parseType(ft.Elt)
		t.isSlice = true
		t.arrayLength = ln

		return t
	default:
		log.Fatalf("internal error: unexpected type: %T", typ)
	}

	return Type{}
}

func (t Type) String() string {
	var s strings.Builder

	if t.name != "" {
		s.WriteString(t.name + " ")
	}
	if t.isVariadic {
		s.WriteString("...")
	}
	if t.isSlice {
		s.WriteString("[")
		if t.arrayLength > 0 {
			s.WriteString(strconv.Itoa(t.arrayLength))
		}
		s.WriteString("]")
	}
	if t.isPointer {
		s.WriteString("*")
	}
	if t.value != "" {
		s.WriteString(t.value)
	}

	return s.String()
}
