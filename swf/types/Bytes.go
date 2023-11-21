package types

import "io"

type Bytes []byte

func (b *Bytes) SWFRead(r DataReader, ctx ReaderContext) (err error) {
	*b, err = io.ReadAll(r)
	if err != nil {
		return err
	}

	return nil
}
