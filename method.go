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

func isContextType(t types.Type) bool {
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	return named.Obj().Name() == "Context"
}
