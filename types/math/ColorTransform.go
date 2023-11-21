package math

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"math"
)

type ColorTransform struct {
	Multiply struct {
		Red   types.Fixed8
		Green types.Fixed8
		Blue  types.Fixed8
		Alpha types.Fixed8
	}
	Add struct {
		Red   int16
		Green int16
		Blue  int16
		Alpha int16
	}
}

func IdentityColorTransform() ColorTransform {
	return ColorTransform{
		Multiply: struct {
			Red   types.Fixed8
			Green types.Fixed8
			Blue  types.Fixed8
			Alpha types.Fixed8
		}{
			Red:   256,
			Green: 256,
			Blue:  256,
			Alpha: 256,
		},
		Add: struct {
			Red   int16
			Green int16
			Blue  int16
			Alpha int16
		}{},
	}
}

func (t ColorTransform) IsIdentity() bool {
	return t.Add.Red == 0 && t.Add.Green == 0 && t.Add.Blue == 0 && t.Add.Alpha == 0 && t.Multiply.Red == 256 && t.Multiply.Green == 256 && t.Multiply.Blue == 256 && t.Multiply.Alpha == 256
}

func (t ColorTransform) ApplyMultiplyToColor(color Color) Color {
	return Color{
		R:     uint8(max(0, min((int64(t.Multiply.Red)*int64(color.R))/(math.MaxUint8+1), 255))),
		G:     uint8(max(0, min((int64(t.Multiply.Green)*int64(color.G))/(math.MaxUint8+1), 255))),
		B:     uint8(max(0, min((int64(t.Multiply.Blue)*int64(color.B))/(math.MaxUint8+1), 255))),
		Alpha: uint8(max(0, min((int64(t.Multiply.Alpha)*int64(color.Alpha))/(math.MaxUint8+1), 255))),
	}
}

func (t ColorTransform) ApplyAdditionToColor(color Color) Color {
	return Color{
		R:     uint8(max(0, min(int64(t.Add.Red)+int64(color.R), 255))),
		G:     uint8(max(0, min(int64(t.Add.Green)+int64(color.G), 255))),
		B:     uint8(max(0, min(int64(t.Add.Blue)+int64(color.B), 255))),
		Alpha: uint8(max(0, min(int64(t.Add.Alpha)+int64(color.Alpha), 255))),
	}
}

func (t ColorTransform) ApplyToColor(color Color) Color {
	return Color{
		R:     uint8(max(0, min((int64(t.Multiply.Red)*int64(color.R))/(math.MaxUint8+1)+int64(t.Add.Red), 255))),
		G:     uint8(max(0, min((int64(t.Multiply.Green)*int64(color.G))/(math.MaxUint8+1)+int64(t.Add.Green), 255))),
		B:     uint8(max(0, min((int64(t.Multiply.Blue)*int64(color.B))/(math.MaxUint8+1)+int64(t.Add.Blue), 255))),
		Alpha: uint8(max(0, min((int64(t.Multiply.Alpha)*int64(color.Alpha))/(math.MaxUint8+1)+int64(t.Add.Alpha), 255))),
	}
}

func (t ColorTransform) Combine(o ColorTransform) ColorTransform {
	return ColorTransform{
		Multiply: struct {
			Red   types.Fixed8
			Green types.Fixed8
			Blue  types.Fixed8
			Alpha types.Fixed8
		}{
			//TODO: maybe needs more than just /(math.MaxUint8+1)
			Red:   types.Fixed8(max(math.MinInt16, min((int64(t.Multiply.Red)*int64(o.Multiply.Red))/(math.MaxUint8+1), math.MaxInt16))),
			Green: types.Fixed8(max(math.MinInt16, min((int64(t.Multiply.Green)*int64(o.Multiply.Green))/(math.MaxUint8+1), math.MaxInt16))),
			Blue:  types.Fixed8(max(math.MinInt16, min((int64(t.Multiply.Blue)*int64(o.Multiply.Blue))/(math.MaxUint8+1), math.MaxInt16))),
			Alpha: types.Fixed8(max(math.MinInt16, min((int64(t.Multiply.Alpha)*int64(o.Multiply.Alpha))/(math.MaxUint8+1), math.MaxInt16))),
		},

		Add: struct {
			Red   int16
			Green int16
			Blue  int16
			Alpha int16
		}{
			Red:   int16(max(math.MinInt16, min(int64(t.Add.Red)+int64(o.Add.Red), math.MaxInt16))),
			Green: int16(max(math.MinInt16, min(int64(t.Add.Green)+int64(o.Add.Green), math.MaxInt16))),
			Blue:  int16(max(math.MinInt16, min(int64(t.Add.Blue)+int64(o.Add.Blue), math.MaxInt16))),
			Alpha: int16(max(math.MinInt16, min(int64(t.Add.Alpha)+int64(o.Add.Alpha), math.MaxInt16))),
		},
	}
}

func ColorTransformFromSWFAlpha(cx types.CXFORMWITHALPHA) (t ColorTransform) {
	t.Multiply.Red = cx.Multiply.Red
	t.Multiply.Green = cx.Multiply.Red
	t.Multiply.Blue = cx.Multiply.Blue
	t.Multiply.Alpha = cx.Multiply.Alpha
	t.Add.Red = cx.Add.Red
	t.Add.Green = cx.Add.Red
	t.Add.Blue = cx.Add.Blue
	t.Add.Alpha = cx.Add.Alpha
	return t
}

func ColorTransformFromSWF(cx types.CXFORM) (t ColorTransform) {
	t.Multiply.Red = cx.Multiply.Red
	t.Multiply.Green = cx.Multiply.Red
	t.Multiply.Blue = cx.Multiply.Blue
	t.Multiply.Alpha = 256
	t.Add.Red = cx.Add.Red
	t.Add.Green = cx.Add.Red
	t.Add.Blue = cx.Add.Blue
	t.Add.Alpha = 0
	return t
}
