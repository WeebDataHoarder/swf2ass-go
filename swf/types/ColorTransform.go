package types

type ColorTransform struct {
	_    struct{} `swfFlags:"root,alignend"`
	Flag struct {
		HasAddTerms  bool
		HasMultTerms bool
	}
	NBits    uint8 `swfBits:",4"`
	Multiply struct {
		Red   Fixed8 `swfBits:"NBits"`
		Green Fixed8 `swfBits:"NBits"`
		Blue  Fixed8 `swfBits:"NBits"`
	} `swfCondition:"Flag.HasMultTerms"`
	Add struct {
		Red   Fixed8 `swfBits:"NBits"`
		Green Fixed8 `swfBits:"NBits"`
		Blue  Fixed8 `swfBits:"NBits"`
	} `swfCondition:"Flag.HasAddTerms"`
}

func (cf *ColorTransform) SWFDefault(ctx ReaderContext) {
	*cf = ColorTransform{}
	cf.Multiply.Red = 256
	cf.Multiply.Green = 256
	cf.Multiply.Blue = 256
}
