package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type DefineFont struct {
	_                struct{} `swfFlags:"root"`
	FontId           uint16
	NumGlyphsEntries uint16
	OffsetTable      []uint16         `swfCount:"TableLength()"`
	ShapeTable       []subtypes.SHAPE `swfCount:"TableLength()"`
}

func (t *DefineFont) TableLength(ctx types.ReaderContext) uint64 {
	return uint64(t.NumGlyphsEntries / 2)
}

func (t *DefineFont) Code() Code {
	return RecordDefineFont
}
