package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type DefineBitsJPEG3 struct {
	_               struct{} `swfFlags:"root"`
	CharacterId     uint16
	AlphaDataOffset uint32
	ImageData       []byte `swfCount:"AlphaDataOffset"`
	BitmapAlphaData types.Bytes
}

func (t *DefineBitsJPEG3) Code() Code {
	return RecordDefineBitsJPEG3
}
