package types

import (
	"strconv"
	"strings"
)

type Depth []uint16

// GetPackedLayer
// Segment depth into specific layers, leaving 2^16 for first, at least 2^8 for second (and 2^8 for third), if no third it'll use whole for second TODO: handle higher depths gracefully
// It is known layers CAN overlap, TODO: check if limiting range might make sense?
// TODO: change this to a truly dynamic mode. might need 2-pass to check for hole overlap
func (d Depth) GetPackedLayer() (layer uint32) {
	if len(d) == 0 {
		return 0
	}
	layer = uint32(d[0]) << 16
	if len(d) > 2 {
		layer |= uint32(d[1]&0xFF) << 8
		layer |= uint32(d[2] & 0xFF)
	} else if len(d) > 1 {
		layer |= uint32(d[1])
	}
	return layer
}

func (d Depth) Equals(o Depth) bool {
	if len(d) != len(o) {
		return false
	}
	for i := range d {
		if d[i] != o[i] {
			return false
		}
	}
	return true
}

func (d Depth) String() string {
	o := make([]string, 0, len(d))
	for _, e := range d {
		o = append(o, strconv.FormatUint(uint64(e), 10))
	}
	return strings.Join(o, ".")
}
