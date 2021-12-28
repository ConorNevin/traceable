package example

import (
	"context"
)

//go:generate ../bin/traceable -type Searcher -output traced_searcher.go

type Searcher interface {
	Search(context.Context, string) error
}
