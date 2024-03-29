package traceable

import (
	"go/types"
	"strconv"
)

type Method struct {
	name       string
	args       []types.Type
	returns    []types.Type
	isVariadic bool
}

func (m Method) acceptsContext() bool {
	if len(m.args) == 0 {
		return false
	}

	for _, a := range m.args {
		if isContextType(a) {
			return true
		}
	}

	return false
}

func (m Method) contextArg() string {
	for i, a := range m.args {
		if isContextType(a) {
			return "a" + strconv.Itoa(i)
		}
	}

	return ""
}

func (m Method) imports() map[string]struct{} {
	imports := make(map[string]struct{})
	for _, t := range m.args {
		for ip := range importsOf(t) {
			imports[ip] = struct{}{}
		}
	}
	for _, t := range m.returns {
		for ip := range importsOf(t) {
			imports[ip] = struct{}{}
		}
	}

	return imports
}

func importsOf(t types.Type) map[string]struct{} {
	switch u := t.(type) {
	case *types.Pointer:
		return importsOf(u.Elem())
	case *types.Map:
		return mergeMaps(importsOf(u.Key()), importsOf(u.Elem()))
	case *types.Array:
		return importsOf(u.Elem())
	case *types.Slice:
		return importsOf(u.Elem())
	case *types.Chan:
		return importsOf(u.Elem())
	case *types.Signature:
		imports := make(map[string]struct{})
		for i := 0; i < u.Params().Len(); i++ {
			imports = mergeMaps(imports, importsOf(u.Params().At(i).Type()))
		}
		for i := 0; i < u.Results().Len(); i++ {
			imports = mergeMaps(imports, importsOf(u.Results().At(i).Type()))
		}

		return imports
	case *types.Named:
		if pkg := u.Obj().Pkg(); pkg != nil {
			return map[string]struct{}{
				pkg.Path(): struct{}{},
			}
		}
		return nil
	default:
		return nil
	}
}

func isContextType(t types.Type) bool {
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	return named.Obj().Name() == "Context"
}

func mergeMaps(a map[string]struct{}, maps ...map[string]struct{}) map[string]struct{} {
	for _, b := range maps {
		for bk, bv := range b {
			if _, ok := a[bk]; !ok {
				a[bk] = bv
			}
		}
	}
	return a
}
