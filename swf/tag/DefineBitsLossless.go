package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type DefineBitsLossless struct {
	_              struct{} `swfFlags:"root"`
	CharacterId    uint16
	Format         uint8
	Width, Height  uint16
	ColorTableSize uint8 `swfCondition:"HasColorTableSize()"`
	ZlibData       types.Bytes
}

func (t *DefineBitsLossless) HasColorTableSize(ctx types.ReaderContext) bool {
	return t.Format == 3
}

func (t *DefineBitsLossless) Code() Code {
	return RecordDefineBitsLossless
}
