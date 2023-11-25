package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type DefineBits struct {
	_           struct{} `swfFlags:"root"`
	CharacterId uint16
	Data        types.UntilEndBytes
}

func (t *DefineBits) Code() Code {
	return RecordDefineBits
}
