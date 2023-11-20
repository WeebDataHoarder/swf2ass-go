package math

import (
	"golang.org/x/exp/constraints"
)

func Lerp[T constraints.Integer | constraints.Float](start, end T, ratio float64) T {
	return start + T(float64(end)-float64(start)*ratio)
}

func LerpVector2[T ~int64 | ~float64](start, end Vector2[T], ratio float64) Vector2[T] {
	return Vector2ToType[float64, T](start.AddVector(end.SubVector(start)).Float64().Multiply(ratio))
}

func LerpColor(start, end Color, ratio float64) Color {
	return Color{
		R:     Lerp(start.R, end.R, ratio),
		G:     Lerp(start.G, end.G, ratio),
		B:     Lerp(start.B, end.B, ratio),
		Alpha: Lerp(start.Alpha, end.Alpha, ratio),
	}
}
