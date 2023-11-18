package types

import (
	"github.com/x448/float16"
)

// Float16 TODO: check if proper values
type Float16 float16.Float16

func (f *Float16) SWFRead(r DataReader, swfVersion uint8) (err error) {
	panic("todo")
}
