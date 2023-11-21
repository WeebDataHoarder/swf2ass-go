package types

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type RenderedObject struct {
	Depth    Depth
	ObjectId uint16

	DrawPathList shapes.DrawPathList

	Clip *ClipPath

	ColorTransform  math.ColorTransform
	MatrixTransform math.MatrixTransform
}

func (o *RenderedObject) GetDepth() Depth {
	if len(o.Depth) > 0 && o.Depth[0] == 0 {
		return o.Depth[1:]
	}
	return o.Depth
}
