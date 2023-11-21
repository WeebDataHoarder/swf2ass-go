package tag

import "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"

type DefineBitsJPEG2 struct {
	_           struct{} `swfFlags:"root"`
	CharacterId uint16
	Data        types.Bytes
}

func (t *DefineBitsJPEG2) Code() Code {
	return RecordDefineBitsJPEG2
}
