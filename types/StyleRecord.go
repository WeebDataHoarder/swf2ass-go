package types

import (
	swfsubtypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	swftypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type StyleRecord interface {
	ApplyColorTransform(transform math.ColorTransform) StyleRecord
}

type LineStyleRecord struct {
	Width swftypes.Twip
	Color math.Color
}

func (r *LineStyleRecord) ApplyColorTransform(transform math.ColorTransform) StyleRecord {
	return &LineStyleRecord{
		Width: r.Width,
		Color: transform.ApplyToColor(r.Color),
	}
}

type FillStyleRecord struct {
	// Fill can be a Color or Gradient
	Fill   any
	Border *LineStyleRecord
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
	}
}

func FillStyleRecordFromSWFFILLSTYLE(fillStyle swfsubtypes.FILLSTYLE) (r *FillStyleRecord) {
	switch fillStyle.FillStyleType {
	case swfsubtypes.FillStyleSolid:
		return &FillStyleRecord{
			Fill: math.Color{
				R:     fillStyle.Color.R(),
				G:     fillStyle.Color.G(),
				B:     fillStyle.Color.B(),
				Alpha: fillStyle.Color.A(),
			},
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
