package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type DefineMorphShape2 struct {
	_                              struct{} `swfFlags:"root"`
	CharacterId                    uint16
	StartBounds, EndBounds         types.Rectangle
	StartEdgeBounds, EndEdgeBounds types.Rectangle
	Reserved                       uint8 `swfBits:",6"`
	UsesNonScalingStrokes          bool
	UsesScalingStrokes             bool
	Offset                         uint32
	MorphFillStyles                subtypes.MORPHFILLSTYLEARRAY
	MorphLineStyles                subtypes.MORPHLINESTYLEARRAY
	StartEdges, EndEdges           subtypes.SHAPE
}

func (t *DefineMorphShape2) Code() Code {
	return RecordDefineMorphShape2
}
