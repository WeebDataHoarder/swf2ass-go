package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type DefineFontAlignZones struct {
	_            struct{} `swfFlags:"root"`
	FontId       uint16
	CSMTableHint uint8 `swfBits:",2"`
	Reserved     uint8 `swfBits:",6"`
	ZoneTable    types.UntilEnd[subtypes.ZONERECORD]
}

func (t *DefineFontAlignZones) Code() Code {
	return RecordDefineFontAlignZones
}
