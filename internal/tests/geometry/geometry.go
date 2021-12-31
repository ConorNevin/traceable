package geometry

import (
	"context"
	"fmt"
	"math"
)

//go:generate ../../../bin/traceable -types Geometry -output geometry_traced.go

type Geometry interface {
	Area(context.Context) (float64, error)
}

type Circle struct {
	radius float64
}

func (c *Circle) Area(_ context.Context) (float64, error) {
	if c.radius < 0 {
		return 0, fmt.Errorf("unable to calculate area for circle with size %.2f", c.radius)
	}

	return math.Pi * c.radius * c.radius, nil
}
