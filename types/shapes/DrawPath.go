package shapes

import "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"

type DrawPath struct {
	Style    StyleRecord
	Clip     *ClipPath
	Commands *Shape
}

func (p DrawPath) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) (r DrawPath) {
	if p.Clip == nil {
		return DrawPath{
			Style:    p.Style,
			Commands: p.Commands.ApplyMatrixTransform(transform, applyTranslation),
		}
	}
	return DrawPath{
		Style:    p.Style,
		Commands: p.Commands.ApplyMatrixTransform(transform, applyTranslation),
		Clip:     p.Clip.ApplyMatrixTransform(transform, applyTranslation),
	}
}

func DrawPathFill(record *FillStyleRecord, shape *Shape, clip *ClipPath) DrawPath {
	return DrawPath{
		Style:    record,
		Commands: shape,
		Clip:     clip,
	}
}

func DrawPathStroke(record *LineStyleRecord, shape *Shape, clip *ClipPath) DrawPath {
	return DrawPath{
		Style:    record,
		Commands: shape,
		Clip:     clip,
	}
}
