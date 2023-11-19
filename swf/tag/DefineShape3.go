package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type DefineShape3 struct {
	_           struct{} `swfFlags:"root,align"`
	ShapeId     uint16
	ShapeBounds types.Rectangle
	Shapes      subtypes.SHAPEWITHSTYLE `swfFlags:"Shape3"`
}

func (t *DefineShape3) Code() Code {
	return RecordDefineShape3
}
