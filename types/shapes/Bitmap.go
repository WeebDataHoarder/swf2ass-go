package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	math2 "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type Bitmap struct {
	List DrawPathList

	Transform math2.MatrixTransform
}

func (b Bitmap) ApplyColorTransform(transform math2.ColorTransform) Fillable {
	b2 := b
	b2.List = b.List.ApplyColorTransform(transform).(DrawPathList)
	b2.Transform = b.Transform
	return b2
}

func (b Bitmap) Fill(shape *Shape) DrawPathList {
	return b.List.Fill(shape)
}

func BitmapFillFromSWF(l DrawPathList, transform types.MATRIX) Bitmap {
	// shape is already in pixel world, but matrix comes as twip
	baseScale := math2.ScaleTransform(math2.NewVector2[float64](1./types.TwipFactor, 1./types.TwipFactor))
	return Bitmap{
		List:      l,
		Transform: math2.MatrixTransformFromSWF(transform).Multiply(baseScale),
	}
}
