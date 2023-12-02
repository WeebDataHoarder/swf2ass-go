package swf

import (
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"github.com/icza/bitio"
	"github.com/ulikunitz/xz/lzma"
	"io"
	"reflect"
)

type Reader struct {
	r      *bitio.Reader
	header types.Header
}

func NewReader(reader io.Reader) (*Reader, error) {
	r := &Reader{}

	var headerData [8]byte
	if _, err := reader.Read(headerData[:]); err != nil {
		return nil, err
	}

	copy(r.header.Signature[:], headerData[:])
	r.header.Version = headerData[3]
	r.header.FileLength = binary.LittleEndian.Uint32(headerData[4:])

	switch r.header.Signature {
	case types.SignatureUncompressed:
		r.r = bitio.NewReader(reader)
	case types.SignatureCompressedZLIB:
		if r.header.Version < 6 {
			return nil, fmt.Errorf("unsupported signature %s", string(r.header.Signature[:]))
		}
		zlibReader, err := zlib.NewReader(reader)
		if err != nil {
			return nil, err
		}
		r.r = bitio.NewReader(zlibReader)
	case types.SignatureCompressedLZMA:
		if r.header.Version < 13 {
			return nil, fmt.Errorf("unsupported signature %s", string(r.header.Signature[:]))
		}
		lzmaReader, err := lzma.NewReader(reader)
		if err != nil {
			return nil, err
		}
		r.r = bitio.NewReader(lzmaReader)
	default:
		return nil, fmt.Errorf("unsupported signature %s", string(r.header.Signature[:]))
	}

	err := types.ReadType(r.r, types.ReaderContext{
		Version: r.header.Version,
	}, &r.header.FrameSize)
	if err != nil {
		return nil, err
	}

	err = types.ReadSI16(r.r, &r.header.FrameRate)
	if err != nil {
		return nil, err
	}

	err = types.ReadU16[uint16](r.r, &r.header.FrameCount)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Reader) Header() types.Header {
	return r.header
}

func (r *Reader) Record() (record *tag.Record, err error) {
	record = &tag.Record{}
	err = types.ReadType(r.r, types.ReaderContext{
		Version: r.header.Version,
	}, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *Reader) Tag() (readTag tag.Tag, err error) {
	record, err := r.Record()
	if err != nil {
		return nil, err
	}

	readTag, err = record.Decode()

	if types.DoParserDebug {
		if readTag == nil {
			fmt.Printf("%d: len %d UNKNOWN\n", record.Code(), len(record.Data))
		} else {
			fmt.Printf("%d: len %d KNOWN %s\n", record.Code(), len(record.Data), reflect.ValueOf(readTag).Elem().Type().Name())
		}
	}

	return readTag, err
}
