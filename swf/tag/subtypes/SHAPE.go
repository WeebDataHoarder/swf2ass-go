package subtypes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"reflect"
)

type SHAPE struct {
	_        struct{} `swfFlags:"root"`
	FillBits uint8    `swfBits:",4"`
	LineBits uint8    `swfBits:",4"`
	Records  SHAPERECORDS
}

type SHAPEWITHSTYLE struct {
	_          struct{} `swfFlags:"root"`
	FillStyles FILLSTYLEARRAY
	LineStyles LINESTYLEARRAY
	FillBits   uint8 `swfBits:",4"`
	LineBits   uint8 `swfBits:",4"`
	Records    SHAPERECORDS
}

type SHAPERECORDS []SHAPERECORD

func (records *SHAPERECORDS) SWFRead(r types.DataReader, ctx types.ReaderContext) (err error) {
	fillBits := uint8(ctx.Root.FieldByName("FillBits").Uint())
	lineBits := uint8(ctx.Root.FieldByName("LineBits").Uint())
	for {
		isEdge, err := types.ReadBool(r)
		if err != nil {
			return err
		}

		if !isEdge {
			rec := StyleChangeRecord{}

			rec.FillBits = fillBits
			rec.LineBits = lineBits

			err = types.ReadType(r, types.ReaderContext{
				Version: ctx.Version,
				Root:    reflect.ValueOf(rec),
				Flags:   ctx.Flags,
			}, &rec)
			if err != nil {
				return err
			}

			if !rec.Flag.NewStyles && !rec.Flag.LineStyle && !rec.Flag.FillStyle1 && !rec.Flag.FillStyle0 && !rec.Flag.MoveTo {
				//end record
				*records = append(*records, &EndShapeRecord{})
				break
			}

			//store new value
			fillBits = rec.FillBits
			lineBits = rec.LineBits

			*records = append(*records, rec)
		} else {
			isStraight, err := types.ReadBool(r)
			if err != nil {
				return err
			}
			if isStraight {
				rec := StraightEdgeRecord{}
				err = types.ReadType(r, types.ReaderContext{
					Version: ctx.Version,
					Root:    reflect.ValueOf(rec),
					Flags:   ctx.Flags,
				}, &rec)
				if err != nil {
					return err
				}
				*records = append(*records, rec)
			} else {
				rec := CurvedEdgeRecord{}

				err = types.ReadType(r, types.ReaderContext{
					Version: ctx.Version,
					Root:    reflect.ValueOf(rec),
					Flags:   ctx.Flags,
				}, &rec)
				if err != nil {
					return err
				}
				*records = append(*records, rec)
			}
		}
	}

	r.Align()

	return nil
}

type EndShapeRecord struct {
}

type StyleChangeRecord struct {
	_    struct{} `swfFlags:"root"`
	Flag struct {
		NewStyles  bool
		LineStyle  bool
		FillStyle1 bool
		FillStyle0 bool
		MoveTo     bool
	}

	MoveBits               uint8          `swfBits:",5" swfCondition:"Flag.MoveTo"`
	MoveDeltaX, MoveDeltaY types.Twip     `swfBits:"MoveBits,signed" swfCondition:"Flag.MoveTo"`
	FillStyle0             uint16         `swfBits:"FillBits" swfCondition:"Flag.FillStyle0"`
	FillStyle1             uint16         `swfBits:"FillBits" swfCondition:"Flag.FillStyle1"`
	LineStyle              uint16         `swfBits:"LineBits" swfCondition:"Flag.LineStyle"`
	FillStyles             FILLSTYLEARRAY `swfFlags:"align" swfCondition:"Flag.NewStyles"`
	LineStyles             LINESTYLEARRAY `swfCondition:"Flag.NewStyles"`

	FillBits uint8 `swfBits:",4" swfCondition:"Flag.NewStyles"`
	LineBits uint8 `swfBits:",4" swfCondition:"Flag.NewStyles"`
}

type StraightEdgeRecord struct {
	_ struct{} `swfFlags:"root"`

	NumBits uint8 `swfBits:",4"`

	GeneralLine bool
	VertLine    bool `swfCondition:"HasVertLine()"`

	DeltaX types.Twip `swfBits:"NumBits+2,signed" swfCondition:"HasDeltaX()"`
	DeltaY types.Twip `swfBits:"NumBits+2,signed" swfCondition:"HasDeltaY()"`
}

func (s *StraightEdgeRecord) HasVertLine(ctx types.ReaderContext) bool {
	return !s.GeneralLine
}

func (s *StraightEdgeRecord) HasDeltaX(ctx types.ReaderContext) bool {
	return s.GeneralLine || !s.VertLine
}

func (s *StraightEdgeRecord) HasDeltaY(ctx types.ReaderContext) bool {
	return s.GeneralLine || s.VertLine
}

type CurvedEdgeRecord struct {
	_ struct{} `swfFlags:"root"`

	NumBits uint8 `swfBits:",4"`

	ControlDeltaX, ControlDeltaY types.Twip `swfBits:"NumBits+2,signed"`
	AnchorDeltaX, AnchorDeltaY   types.Twip `swfBits:"NumBits+2,signed"`
}

type SHAPERECORD interface {
}
