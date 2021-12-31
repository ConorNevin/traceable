package example

import (
	"context"
	"io"
	"net/http"
)

//go:generate ../bin/traceable -type Embedded -output embedded_types_traced.go

type Base interface {
	FunctionOne(context.Context, func(context.Context, io.Reader) error) error
	FunctionThree(context.Context, []http.Request) error
}

type Embedded interface {
	Base
	FunctionOne(context.Context, func(context.Context, io.Reader) error) error
	FunctionTwo(context.Context, io.Writer) error
}
