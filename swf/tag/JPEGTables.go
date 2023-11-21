package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type JPEGTables struct {
	_    struct{} `swfFlags:"root"`
	Data types.Bytes
}

func (t *JPEGTables) Code() Code {
	return RecordJPEGTables
}
