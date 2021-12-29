package example

import (
	"context"
)

//go:generate ../bin/traceable -type Searcher -output searcher_traced.go

type Searcher interface {
	Search(context.Context, string) error
	SearchAll(context.Context, ...string) (chan<- string, error)
	StoreAll(context.Context, <-chan string) error
}
