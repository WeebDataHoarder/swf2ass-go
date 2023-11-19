package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type DefineShape4 struct {
	_                     struct{} `swfFlags:"root,align"`
	ShapeId               uint16
	ShapeBounds           types.Rectangle
	EdgeBounds            types.Rectangle
	Reserved              uint8 `swfBits:",5"`
	UsesFillWindingRule   bool
	UsesNonScalingStrokes bool
	UsesScalingStrokes    bool
	Shapes                subtypes.SHAPEWITHSTYLE `swfFlags:"Shape4"`
}

func (t *DefineShape4) Code() Code {
	return RecordDefineShape4
}