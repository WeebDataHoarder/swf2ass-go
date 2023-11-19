package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type DefineShape struct {
	_           struct{} `swfFlags:"root"`
	ShapeId     uint16
	ShapeBounds types.Rectangle
	Shapes      subtypes.SHAPEWITHSTYLE
}

func (t *DefineShape) Code() Code {
	return RecordDefineShape
}
