package types

import (
	"fmt"
	"math"
)

type Color struct {
	R, G, B, Alpha uint8
}

func (c *Color) Equals(o *Color, alpha bool) bool {
	if !alpha {
		return c.R == o.R && c.G == o.G && c.B == o.B
	}
	return *c == *o
}

func intPowD(a, b uint8) int64 {
	return (int64(a) - int64(b)) * (int64(a) - int64(b))
}

func (c *Color) Distance(o *Color, alpha bool) float64 {
	dist := intPowD(c.R, o.R) + intPowD(c.G, o.G) + intPowD(c.B, o.B)
	if alpha {
		dist += intPowD(c.Alpha, o.Alpha)
	}
	return math.Sqrt(float64(dist))
}

func (c *Color) String() string {
	return fmt.Sprintf("rgba(%d,%d,%d,%d)", c.R, c.G, c.B, c.Alpha)
}
