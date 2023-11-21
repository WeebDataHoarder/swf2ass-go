package types

import (
	"bytes"
	"github.com/icza/bitio"
	"testing"
)

func TestReadSB(t *testing.T) {
	val := int8(-7)
	data := []byte{uint8(val)}
	r := bitio.NewReader(bytes.NewReader(data))

	result, err := ReadSB[int64](r, 8)
	if err != nil {
		t.Fatal(err)
	}
	if result != int64(val) {
		t.Fatal("does not match")
	}

}
