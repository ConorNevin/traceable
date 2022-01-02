package embedded_interface

import (
	"context"
	"io"
	"net/http"

	"github.com/ConorNevin/traceable/internal/tests/embedded_interface/nested"
)

//go:generate ../../../bin/traceable -types Embedded -output embedded_types_traced.go
//go:generate ../../../bin/traceable -types AnotherEmbedded -output another_embedded_types_traced.go

type Base interface {
	FunctionOne(context.Context, func(context.Context, io.Reader) error) error
	FunctionThree(context.Context, []http.Request) error
}

type Embedded interface {
	Base
	FunctionOne(context.Context, func(context.Context, io.Reader) error) error
	FunctionTwo(context.Context, io.Writer) error
}

type AnotherEmbedded interface {
	nested.Faux
	FauxDu(context.Context) (string, func() error, error)
}
