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
	b2.List = b2.List.ApplyColorTransform(transform).(DrawPathList)
	return b2
}

func (b Bitmap) ApplyMatrixTransform(transform math2.MatrixTransform, applyTranslation bool) Fillable {
	b2 := b
	if !applyTranslation {
		panic("not supported")
	}
	b2.Transform = transform.Combine(b2.Transform)
	return b2
}

func (b Bitmap) Fill(shape Shape) DrawPathList {
	return b.List.ApplyMatrixTransform(b.Transform, true).Fill(shape)
}

func BitmapFillFromSWF(l DrawPathList, transform types.MATRIX) Bitmap {
	return Bitmap{
		List: l,
		// shape is already in pixel world, but matrix comes as twip
		Transform: math2.MatrixTransformFromSWF(transform, 1./types.TwipFactor),
	}
}
