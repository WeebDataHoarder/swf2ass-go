package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"io"
)

type DefineBitsJPEG4 struct {
	_               struct{} `swfFlags:"root"`
	CharacterId     uint16
	AlphaDataOffset uint32
	DeblockParam    types.Fixed8
	ImageData       []byte
	BitmapAlphaData []byte
}

func (t *DefineBitsJPEG4) SWFRead(r types.DataReader, ctx types.ReaderContext) (err error) {
	err = types.ReadU16(r, &t.CharacterId)
	if err != nil {
		return err
	}
	err = types.ReadU32(r, &t.AlphaDataOffset)
	if err != nil {
		return err
	}
	err = types.ReadU16(r, &t.DeblockParam)
	if err != nil {
		return err
	}

	t.ImageData = make([]byte, t.AlphaDataOffset)
	_, err = io.ReadFull(r, t.ImageData)
	if err != nil {
		return err
	}

	t.BitmapAlphaData, err = io.ReadAll(r)
	if err != nil {
		return err
	}

	return nil
}

func (t *DefineBitsJPEG4) Code() Code {
	return RecordDefineBitsJPEG4
}