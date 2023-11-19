package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"io"
)

type JPEGTables struct {
	_    struct{} `swfFlags:"root"`
	Data []byte
}

func (t *JPEGTables) SWFRead(r types.DataReader, ctx types.ReaderContext) (err error) {
	t.Data, err = io.ReadAll(r)
	if err != nil {
		return err
	}

	return nil
}

func (t *JPEGTables) Code() Code {
	return RecordJPEGTables
}
