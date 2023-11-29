package shapes

import "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"

type DrawPath struct {
	Style StyleRecord
	Shape Shape
}

func (p DrawPath) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) (r DrawPath) {
	return DrawPath{
		Style: p.Style.ApplyMatrixTransform(transform, applyTranslation),
		Shape: p.Shape.ApplyMatrixTransform(transform, applyTranslation),
	}
}

func (p DrawPath) ApplyColorTransform(transform math.ColorTransform) (r DrawPath) {
	return DrawPath{
		Style: p.Style.ApplyColorTransform(transform),
		Shape: p.Shape,
	}
}

func DrawPathFill(record *FillStyleRecord, shape Shape) DrawPath {
	return DrawPath{
		Style: record,
		Shape: shape,
	}
}

func DrawPathStroke(record *LineStyleRecord, shape Shape) DrawPath {
	return DrawPath{
		Style: record,
		Shape: shape,
	}
}
