package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type DefineText2 struct {
	_           struct{} `swfFlags:"root"`
	CharacterId uint16
	Bounds      types.RECT
	Matrix      types.MATRIX
	GlyphBits   uint8
	AdvanceBits uint8
	TextRecords subtypes.TEXTRECORDS `swfFlags:"Text2"`
}

func (t *DefineText2) Code() Code {
	return RecordDefineText2
}
