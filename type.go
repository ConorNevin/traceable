package traceable

import (
	"fmt"
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
	packageName   string
	value         string
	isPointer     bool
	isSlice       bool
	isVariadic    bool
	isFunc        bool
	funcArguments []*Parameter
	funcReturns   []*Parameter

	isChan  bool
	chanDir ChanDir

	isMap            bool
	mapKey, mapValue *Type

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
	if t.isFunc {
		args := make([]string, len(t.funcArguments))
		for i, p := range t.funcArguments {
			args[i] = p.typ.String()
		}

		rets := make([]string, len(t.funcReturns))
		for i, p := range t.funcReturns {
			rets[i] = p.typ.String()
		}
		retStr := strings.Join(rets, ", ")
		if len(rets) > 1 {
			retStr = "(" + retStr + ")"
		}

		s.WriteString("func(" + strings.Join(args, ",") + ") " + retStr)
		return s.String()
	}
	if t.isMap {
		s.WriteString(fmt.Sprintf("map[%s]%s", t.mapKey.String(), t.mapValue.String()))
		return s.String()
	}
	if t.isChan {
		switch t.chanDir {
		case ChanSend:
			s.WriteString("chan<- ")
		case ChanRecv:
			s.WriteString("<-chan ")
		default:
			s.WriteString("chan ")
		}
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

type ChanDir uint8

const (
	ChanSend ChanDir = 1
	ChanRecv ChanDir = 2
)
