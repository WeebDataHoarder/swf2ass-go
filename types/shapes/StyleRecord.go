package shapes

import (
	swfsubtypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	swftypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type StyleRecord interface {
	ApplyColorTransform(transform math.ColorTransform) StyleRecord
}

type LineStyleRecord struct {
	Width float64
	Color math.Color
	Blur  float64
}

func (r *LineStyleRecord) ApplyColorTransform(transform math.ColorTransform) StyleRecord {
	return &LineStyleRecord{
		Width: r.Width,
		Color: transform.ApplyToColor(r.Color),
		Blur:  r.Blur,
	}
}

type FillStyleRecord struct {
	// Fill can be a math.Color or Gradient
	Fill   any
	Border *LineStyleRecord
	Blur   float64
}

func (r *FillStyleRecord) ApplyColorTransform(transform math.ColorTransform) StyleRecord {
	fill := r.Fill
	if color, ok := fill.(math.Color); ok {
		fill = transform.ApplyToColor(color)
	} else if gradient, ok := fill.(Gradient); ok {
		fill = gradient.ApplyColorTransform(transform)
	}
	return &FillStyleRecord{
		Border: r.Border,
		Fill:   fill,
		Blur:   r.Blur,
	}
}

func FillStyleRecordFromSWF(fillType swfsubtypes.FillStyleType, color swftypes.Color, gradient swfsubtypes.GRADIENT, gradientMatrix swftypes.MATRIX) (r *FillStyleRecord) {
	switch fillType {
	case swfsubtypes.FillStyleSolid:
		return &FillStyleRecord{
			Fill: math.Color{
				R:     color.R(),
				G:     color.G(),
				B:     color.B(),
				Alpha: color.A(),
			},
		}
	case swfsubtypes.FillStyleLinearGradient:
		return &FillStyleRecord{
			Fill: LinearGradientFromSWF(gradient.Records, gradientMatrix, gradient.SpreadMode, gradient.InterpolationMode),
		}
		//TODO other styles
	}

	return &FillStyleRecord{
		Fill: math.Color{
			R:     0,
			G:     0,
			B:     0,
			Alpha: 0,
		},
	}
}

func FillStyleRecordFromSWFMORPHFILLSTYLEStart(fillStyle swfsubtypes.MORPHFILLSTYLE) (r *FillStyleRecord) {
	return FillStyleRecordFromSWF(fillStyle.FillStyleType, fillStyle.StartColor, fillStyle.Gradient.StartGradient(), fillStyle.StartGradientMatrix)
}

func FillStyleRecordFromSWFMORPHFILLSTYLEEnd(fillStyle swfsubtypes.MORPHFILLSTYLE) (r *FillStyleRecord) {
	return FillStyleRecordFromSWF(fillStyle.FillStyleType, fillStyle.EndColor, fillStyle.Gradient.EndGradient(), fillStyle.EndGradientMatrix)
}
