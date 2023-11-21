package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type DefineBitsJPEG4 struct {
	_               struct{} `swfFlags:"root"`
	CharacterId     uint16
	AlphaDataOffset uint32
	DeblockParam    types.Fixed8
	ImageData       []byte `swfCount:"AlphaDataOffset"`
	BitmapAlphaData types.Bytes
}

func (t *DefineBitsJPEG4) Code() Code {
	return RecordDefineBitsJPEG4
}
