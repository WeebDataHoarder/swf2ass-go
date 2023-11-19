package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"io"
)

type DefineBitsJPEG2 struct {
	_           struct{} `swfFlags:"root"`
	CharacterId uint16
	Data        []byte
}

func (t *DefineBitsJPEG2) SWFRead(r types.DataReader, ctx types.ReaderContext) (err error) {
	err = types.ReadU16(r, &t.CharacterId)
	if err != nil {
		return err
	}

	t.Data, err = io.ReadAll(r)
	if err != nil {
		return err
	}

	return nil
}

func (t *DefineBitsJPEG2) Code() Code {
	return RecordDefineBitsJPEG2
}
