package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"reflect"
)

type DefineSprite struct {
	_           struct{} `swfFlags:"root"`
	SpriteId    uint16
	FrameCount  uint16
	ControlTags []Tag
}

func (t *DefineSprite) SWFRead(r types.DataReader, ctx types.ReaderContext) (err error) {
	err = types.ReadU16(r, &t.SpriteId)
	if err != nil {
		return err
	}
	err = types.ReadU16(r, &t.FrameCount)

	for {
		record := &Record{}
		err = types.ReadType(r, types.ReaderContext{
			Version: ctx.Version,
		}, record)
		if err != nil {
			return err
		}

		readTag, err := record.Decode()
		if err != nil {
			return err
		}

		if readTag == nil {
			//not decoded
			continue
		}

		if types.DoParserDebug {
			if readTag == nil {
				fmt.Printf("SPRITE %d: len %d UNKNOWN\n", record.Code(), len(record.Data))
			} else {
				fmt.Printf("SPRITE %d: len %d KNOWN %s\n", record.Code(), len(record.Data), reflect.ValueOf(readTag).Elem().Type().Name())
			}
		}

		t.ControlTags = append(t.ControlTags, readTag)

		if readTag.Code() == RecordEnd {
			break
		}
	}

	return nil
}

func (t *DefineSprite) Code() Code {
	return RecordDefineSprite
}
