package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type DefineText struct {
	_           struct{} `swfFlags:"root"`
	CharacterId uint16
	Bounds      types.RECT
	Matrix      types.MATRIX
	GlyphBits   uint8
	AdvanceBits uint8
	TextRecords subtypes.TEXTRECORDS
}

func (t *DefineText) Code() Code {
	return RecordDefineText
}
