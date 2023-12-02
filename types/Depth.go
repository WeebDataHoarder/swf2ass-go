package types

import (
	"slices"
	"strconv"
	"strings"
)

type Depth []uint16

// GetPackedLayer
// Segment depth into specific layers, leaving 2^16 for first, at least 2^8 for second (and 2^8 for third), if no third it'll use whole for second TODO: handle higher depths gracefully
// It is known layers CAN overlap, TODO: check if limiting range might make sense?
// TODO: change this to a truly dynamic mode. might need 2-pass to check for hole overlap
// libass reads this onto an int32, TODO it overflows onto negative?
// Additionally the order of fields has extra read order. earlier lines have lower effective layer than later lines, if layers are equal
func (d Depth) GetPackedLayer() (layer int32) {
	if len(d) == 0 {
		return 0
	}

	layer = int32(d[0]) << 16
	if len(d) > 2 {
		layer |= int32(d[1]&0xFF) << 8
		layer |= (int32(d[2]&0x7F) << 1) | 0
	} else if len(d) > 1 {
		layer |= (int32(d[1]&0x7FFF) << 1) | 1
	}
	return layer
}

func DepthFromString(layer string) (d Depth, err error) {
	layers := strings.Split(layer, ".")
	d = make(Depth, 0, len(layers))
	for _, l := range layers {
		e, err := strconv.ParseUint(l, 10, 16)
		if err != nil {
			return nil, err
		}
		d = append(d, uint16(e))
	}
	return d, nil
}

func DepthFromPackedLayer(layer int32) (d Depth) {
	d = append(d, uint16(layer>>16))
	if layer&1 == 1 {
		d = append(d, uint16((layer>>1)&0x7FFF))
	} else {
		d = append(d, uint16((layer>>8)&0xFF))
		d = append(d, uint16((layer>>1)&0x7F))
	}
	return d
}

func (d Depth) Compare(o Depth) int {
	return slices.Compare(d, o)
}

func (d Depth) Equals(o Depth) bool {
	return d.Compare(o) == 0
}

func (d Depth) String() string {
	o := make([]string, 0, len(d))
	for _, e := range d {
		o = append(o, strconv.FormatUint(uint64(e), 10))
	}
	return strings.Join(o, ".")
}
