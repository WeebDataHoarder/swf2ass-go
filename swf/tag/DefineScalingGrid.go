package tag

import "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"

type DefineScalingGrid struct {
	CharacterId uint16
	Splitter    types.Rectangle
}

func (t *DefineScalingGrid) Code() Code {
	return RecordDefineScalingGrid
}