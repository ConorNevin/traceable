package subpackage

//go:generate ../../../bin/traceable -types FooBar -output traced/foobar.go

import (
	"context"
)

type FooBar interface {
	Foo(context.Context) error
}
