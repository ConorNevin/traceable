package traceable

import (
	"strconv"
)

type Method struct {
	name    string
	args    []*Parameter
	returns []*Parameter
}

func (m Method) acceptsContext() bool {
	return len(m.args) > 0 && m.args[0].typ.value == "Context"
}

func (m Method) contextArg() string {
	for i, a := range m.args {
		if a.typ.value == "Context" {
			return "a" + strconv.Itoa(i)
		}
	}

	return ""
}
