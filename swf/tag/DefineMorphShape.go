package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type DefineMorphShape struct {
	_                      struct{} `swfFlags:"root"`
	CharacterId            uint16
	StartBounds, EndBounds types.RECT
	Offset                 uint32
	MorphFillStyles        subtypes.MORPHFILLSTYLEARRAY
	MorphLineStyles        subtypes.MORPHLINESTYLEARRAY
	StartEdges, EndEdges   subtypes.SHAPE
}

func (t *DefineMorphShape) Code() Code {
	return RecordDefineMorphShape
}
