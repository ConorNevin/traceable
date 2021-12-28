package traceable

import (
	"strconv"
	"strings"
)

type Parameter struct {
	name string
	typ  *Type
}

func (p *Parameter) String() string {
	return p.typ.String()
}

type Type struct {
	packageName string
	value       string
	isPointer   bool
	isSlice     bool
	isVariadic  bool

	arrayLength int
}

func (t Type) TypeName() string {
	var s string
	if t.packageName != "" {
		s += t.packageName + "."
	}

	return s + t.value
}

func (t Type) String() string {
	var s strings.Builder

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
	if t.packageName != "" {
		s.WriteString(t.packageName + ".")
	}
	if t.value != "" {
		s.WriteString(t.value)
	}

	return s.String()
}
