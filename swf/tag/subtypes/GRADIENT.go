package subtypes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"slices"
)

type GradientSpreadMode uint8

const (
	GradientSpreadPad = GradientSpreadMode(iota)
	GradientSpreadReflect
	GradientSpreadRepeat
	GradientSpreadReserved
)

type GradientInterpolationMode uint8

const (
	GradientInterpolationRGB = GradientInterpolationMode(iota)
	GradientInterpolationLinearRGB
	GradientSpreadReserved2
	GradientSpreadReserved3
)

type GRADIENT struct {
	_                 struct{}                  `swfFlags:"root"`
	SpreadMode        GradientSpreadMode        `swfBits:",2"`
	InterpolationMode GradientInterpolationMode `swfBits:",2"`
	NumGradients      uint8                     `swfBits:",4"`
	Records           []GRADRECORD              `swfCount:"NumGradients"`

	BogusCheck struct{} `swfCondition:"BogusCheckField()"`
}

func (g *GRADIENT) BogusCheckField(ctx types.ReaderContext) bool {
	if g.NumGradients < 1 {
		panic("wrong length")
	}
	if g.SpreadMode != GradientSpreadPad && g.SpreadMode != GradientSpreadReflect && g.SpreadMode != GradientSpreadRepeat {
		panic("wrong spread")
	}
	if g.InterpolationMode != GradientInterpolationRGB && g.InterpolationMode != GradientInterpolationLinearRGB {
		panic("wrong interpolation")
	}
	return false
}

type FOCALGRADIENT struct {
	_                 struct{}                  `swfFlags:"root"`
	SpreadMode        GradientSpreadMode        `swfBits:",2"`
	InterpolationMode GradientInterpolationMode `swfBits:",2"`
	NumGradients      uint8                     `swfBits:",4"`
	Records           []GRADRECORD              `swfCount:"NumGradients"`
	FocalPoint        types.Fixed8

	BogusCheck struct{} `swfCondition:"BogusCheckField()"`
}

func (g *FOCALGRADIENT) BogusCheckField(ctx types.ReaderContext) bool {
	if g.NumGradients < 1 {
		panic("wrong length")
	}
	if g.SpreadMode != GradientSpreadPad && g.SpreadMode != GradientSpreadReflect && g.SpreadMode != GradientSpreadRepeat {
		panic("wrong spread")
	}
	if g.InterpolationMode != GradientInterpolationRGB && g.InterpolationMode != GradientInterpolationLinearRGB {
		panic("wrong interpolation")
	}
	return false
}

type GRADRECORD struct {
	Ratio uint8
	Color types.Color
}

func (g *GRADRECORD) SWFDefault(ctx types.ReaderContext) {
	if slices.Contains(ctx.Flags, "Shape3") || slices.Contains(ctx.Flags, "Shape4") {
		g.Color = &types.RGBA{}
	} else {
		g.Color = &types.RGB{}
	}
}