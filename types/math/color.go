package math

import (
	"fmt"
	"math"
)

type PackedColor uint32

func NewPackedColor(r, g, b, a uint8) PackedColor {
	return PackedColor((uint64(a) << 24) | (uint64(r) << 16) | (uint64(g) << 8) | uint64(b))
}

func (c PackedColor) Alpha() uint8 {
	return uint8((c >> 24) & 0xFF)
}

func (c PackedColor) R() uint8 {
	return uint8((c >> 16) & 0xFF)
}

func (c PackedColor) G() uint8 {
	return uint8((c >> 8) & 0xFF)
}

func (c PackedColor) B() uint8 {
	return uint8(c & 0xFF)
}

func (c PackedColor) Color() Color {
	return Color{
		R:     c.R(),
		G:     c.G(),
		B:     c.B(),
		Alpha: c.Alpha(),
	}
}

type Color struct {
	R, G, B, Alpha uint8
}

func (c Color) Packed() PackedColor {
	return NewPackedColor(c.R, c.G, c.B, c.Alpha)
}

func (c Color) Equals(o Color, alpha bool) bool {
	if !alpha {
		return c.R == o.R && c.G == o.G && c.B == o.B
	}
	return c == o
}

func intPowD(a, b uint8) int64 {
	return (int64(a) - int64(b)) * (int64(a) - int64(b))
}

func (c Color) Distance(o Color, alpha bool) float64 {
	dist := intPowD(c.R, o.R) + intPowD(c.G, o.G) + intPowD(c.B, o.B)
	if alpha {
		dist += intPowD(c.Alpha, o.Alpha)
	}
	return math.Sqrt(float64(dist))
}

func (c Color) ToLinearRGB() Color {
	return Color{
		R:     uint8(math.Pow(float64(c.R)/255, 2.2) * 255),
		G:     uint8(math.Pow(float64(c.G)/255, 2.2) * 255),
		B:     uint8(math.Pow(float64(c.B)/255, 2.2) * 255),
		Alpha: uint8(math.Pow(float64(c.Alpha)/255, 2.2) * 255),
	}
}

func (c Color) ToSRGB() Color {
	return Color{
		R:     uint8(math.Pow(float64(c.R)/255, 0.4545) * 255),
		G:     uint8(math.Pow(float64(c.G)/255, 0.4545) * 255),
		B:     uint8(math.Pow(float64(c.B)/255, 0.4545) * 255),
		Alpha: uint8(math.Pow(float64(c.Alpha)/255, 0.4545) * 255),
	}
}

func (c Color) String() string {
	return fmt.Sprintf("rgba(%d,%d,%d,%d)", c.R, c.G, c.B, c.Alpha)
}
