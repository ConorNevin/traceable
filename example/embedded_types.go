package example

import (
	"context"
	"io"
)

//go:generate ../bin/traceable -type Embedded -output embedded_types_traced.go

type Embedded interface {
	FunctionOne(context.Context, func(context.Context, io.Reader) error) error
	FunctionTwo(context.Context, io.Writer) error
}
