package tag

type Tag interface {
	Code() Code
}

type Record struct {
	_                struct{} `swfFlags:"root,align"`
	TagCodeAndLength uint16
	ExtraLength      uint32 `swfCondition:"HasExtraLength()"`
	Data             []byte `swfCount:"DataLength()"`
}

func (r *Record) HasExtraLength(swfVersion uint8) bool {
	return (r.TagCodeAndLength & 0x3f) == 0x3f
}

func (r *Record) DataLength(swfVersion uint8) uint64 {
	if (r.TagCodeAndLength & 0x3f) == 0x3f {
		return uint64(r.ExtraLength)
	}
	return uint64(r.TagCodeAndLength & 0x3f)
}

func (r *Record) Code() Code {
	return Code(r.TagCodeAndLength >> 6)
}
