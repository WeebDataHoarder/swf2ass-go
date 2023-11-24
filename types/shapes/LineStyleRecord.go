package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type LineStyleRecord struct {
	Width float64
	Color math.Color
	Blur  float64
}

func (r *LineStyleRecord) StrokeWidth(transform math.MatrixTransform) float64 {
	// Flash renders strokes with a 1px minimum width.
	minWidth := transform.MinimumStrokeWidth()
	return 0.5 * max(r.Width, minWidth)
}

func (r *LineStyleRecord) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) StyleRecord {
	return r
}

func (r *LineStyleRecord) ApplyColorTransform(transform math.ColorTransform) StyleRecord {
	return &LineStyleRecord{
		Width: r.Width,
		Color: transform.ApplyToColor(r.Color),
		Blur:  r.Blur,
	}
}
