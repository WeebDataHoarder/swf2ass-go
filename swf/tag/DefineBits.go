package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"io"
)

type DefineBits struct {
	_           struct{} `swfFlags:"root"`
	CharacterId uint16
	Data        []byte
}

func (t *DefineBits) SWFRead(r types.DataReader, ctx types.ReaderContext) (err error) {
	err = types.ReadU16(r, &t.CharacterId)

	t.Data, err = io.ReadAll(r)
	if err != nil {
		return err
	}

	return nil
}

func (t *DefineBits) Code() Code {
	return RecordDefineBits
}
