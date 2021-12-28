package example

import (
	"context"
)

//go:generate ../bin/traceable -type Searcher -output searcher_traced.go

type Searcher interface {
	Search(context.Context, string) error
}
