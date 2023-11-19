package types

type Matrix struct {
	_                        struct{} `swfFlags:"root,alignend"`
	HasScale                 bool
	NScaleBits               uint8 `swfCondition:"HasScale" swfBits:",5"`
	ScaleX, ScaleY           Fixed `swfCondition:"HasScale" swfBits:"NScaleBits"`
	HasRotate                bool
	NRotateBits              uint8 `swfCondition:"HasRotate" swfBits:",5"`
	RotateSkew0, RotateSkew1 Fixed `swfCondition:"HasRotate" swfBits:"NRotateBits"`
	NTranslateBits           uint8 `swfBits:",5"`
	TranslateX, TranslateY   Twip  `swfBits:"NTranslateBits,signed"`
}

func (matrix *Matrix) SWFDefault(ctx ReaderContext) {
	*matrix = Matrix{}
	matrix.ScaleX = 1 << 16
	matrix.ScaleY = 1 << 16
}
