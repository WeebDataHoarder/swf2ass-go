package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type DefineShape2 struct {
	_           struct{} `swfFlags:"root,align"`
	ShapeId     uint16
	ShapeBounds types.Rectangle
	Shapes      subtypes.SHAPEWITHSTYLE
}

func (t *DefineShape2) Code() Code {
	return RecordDefineShape2
}