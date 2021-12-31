package searcher

import (
	"context"
)

//go:generate ../../../bin/traceable -types Searcher -output searcher_traced.go

type Stringer interface {
	String() error
}

type Searcher interface {
	Search(context.Context, string) error
	SearchAll(context.Context, ...string) (chan<- string, error)
	StoreAll(context.Context, <-chan string) error
	StoreMap(context.Context, map[int8]string) error
	StoreInterface(context.Context, Stringer) (int, error)
	One(context.Context, int, int, string) error
}
