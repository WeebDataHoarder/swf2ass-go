package math

import (
	"golang.org/x/exp/constraints"
)

func Lerp[T constraints.Integer | constraints.Float](start, end T, ratio float64) T {
	return start + T((float64(end)-float64(start))*ratio)
}

func LerpVector2[T ~int64 | ~float64](start, end Vector2[T], ratio float64) Vector2[T] {
	return start.AddVector(Vector2ToType[float64, T](end.SubVector(start).Float64().Multiply(ratio)))
}

func LerpColor(start, end Color, ratio float64) Color {
	return Color{
		R:     Lerp(start.R, end.R, ratio),
		G:     Lerp(start.G, end.G, ratio),
		B:     Lerp(start.B, end.B, ratio),
		Alpha: Lerp(start.Alpha, end.Alpha, ratio),
	}
}

func LerpMatrix(start, end MatrixTransform, ratio float64) MatrixTransform {
	// TODO: Lerping a matrix element-wise is geometrically wrong,
	// but I doubt Flash is decomposing the matrix into scale-rotate-translate?

	return NewMatrixTransform(
		LerpVector2(NewVector2(start.GetA(), start.GetD()), NewVector2(end.GetD(), end.GetC()), ratio),
		LerpVector2(NewVector2(start.GetB(), start.GetC()), NewVector2(end.GetB(), end.GetC()), ratio),
		LerpVector2(NewVector2(start.GetTX(), start.GetTY()), NewVector2(end.GetTX(), end.GetTY()), ratio),
	)
}
