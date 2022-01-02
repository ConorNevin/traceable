package nested

import (
	"context"
)

type FauxReturn string

type Faux interface {
	Foo(context.Context) FauxReturn
}
