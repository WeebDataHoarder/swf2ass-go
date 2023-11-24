package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type DrawPathListFill DrawPathList

func (f DrawPathListFill) ApplyColorTransform(transform math.ColorTransform) Fillable {
	return DrawPathListFill(DrawPathList(f).ApplyColorTransform(transform).(DrawPathList))
}

func (f DrawPathListFill) Fill(shape *Shape) DrawPathList {
	return DrawPathList(f)

}

func DrawPathListFillFromSWF(l DrawPathList, transform types.MATRIX) DrawPathListFill {
	// shape is already in pixel world, but matrix comes as twip
	baseScale := math.ScaleTransform(math.NewVector2[float64](1./types.TwipFactor, 1./types.TwipFactor))
	t := math.MatrixTransformFromSWF(transform).Multiply(baseScale)
	return DrawPathListFill(l.ApplyMatrixTransform(t, true))
}
