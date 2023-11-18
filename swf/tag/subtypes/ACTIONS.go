package subtypes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"github.com/icza/bitio"
)

type CLIPACTIONS struct {
	Reserved      uint16
	AllEventFlags CLIPEVENTFLAGS
	Records       CLIPACTIONRECORDS
}

type CLIPACTIONRECORDS []CLIPACTIONRECORD

func (records *CLIPACTIONRECORDS) SWFRead(r types.DataReader, swfVersion uint8) (err error) {
	for {
		var flags CLIPEVENTFLAGS
		err = types.ReadType(r, swfVersion, &flags)
		if err != nil {
			return err
		}
		if flags.IsEnd() {
			break
		}
		record := CLIPACTIONRECORD{
			EventFlags: flags,
		}
		err = types.ReadType(r, swfVersion, &record)
		if err != nil {
			return err
		}
		*records = append(*records, record)
	}
	return nil
}

type CLIPACTIONRECORD struct {
	EventFlags       CLIPEVENTFLAGS
	ActionRecordSize uint32
	KeyCode          uint8
	Actions          []ACTIONRECORD
}

func (clipRecord *CLIPACTIONRECORD) SWFRead(r types.DataReader, swfVersion uint8) (err error) {
	err = types.ReadU32(r, &clipRecord.ActionRecordSize)
	if err != nil {
		return err
	}
	countReader := bitio.NewCountReader(r)
	if clipRecord.EventFlags.KeyPress {
		err = types.ReadU8(r, &clipRecord.KeyCode)
		if err != nil {
			return err
		}
	}

	//TODO: check
	for uint32(countReader.BitsCount/8) < clipRecord.ActionRecordSize {
		var record ACTIONRECORD
		err = types.ReadType(countReader, swfVersion, &record)
		if err != nil {
			return err
		}
		clipRecord.Actions = append(clipRecord.Actions, record)
	}

	return nil
}

type CLIPEVENTFLAGS struct {
	//align?
	_ struct{} `swfFlags:"root"`

	KeyUp      bool
	KeyDown    bool
	MouseUp    bool
	MouseDown  bool
	MouseMove  bool
	Unload     bool
	EnterFrame bool
	Load       bool

	DragOver       bool
	RollOut        bool
	RollOver       bool
	ReleaseOutside bool
	Release        bool
	Press          bool
	Initialize     bool
	EventData      bool

	//SWF 6 or later

	Reserved1 uint8 `swfBits:",5" swfCondition:"IsSWF6OrGreater()"`
	Construct bool  `swfCondition:"IsSWF6OrGreater()"`
	KeyPress  bool  `swfCondition:"IsSWF6OrGreater()"`
	DragOut   bool  `swfCondition:"IsSWF6OrGreater()"`

	Reserved2 uint8 `swfBits:",8" swfCondition:"IsSWF6OrGreater()"`
}

func (f *CLIPEVENTFLAGS) IsEnd() bool {
	return *f == CLIPEVENTFLAGS{}
}

func (f *CLIPEVENTFLAGS) IsSWF6OrGreater(swfVersion uint8) bool {
	return swfVersion >= 6
}

type ACTIONRECORD struct {
	_          struct{} `swfFlags:"root"`
	ActionCode uint8
	Length     uint16  `swfCondition:"HasActionLength()"`
	Data       []uint8 `swfCount:"Length"`
}

func (a *ACTIONRECORD) HasActionLength(swfVersion uint8) bool {
	return a.ActionCode > 0x80
}
