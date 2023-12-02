package processing

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"slices"
	"time"
)

type EventLineHeaders []EventLineHeader

// Sort Sorts each line
func (hs EventLineHeaders) Sort() {
	slices.SortStableFunc(hs, SortEventLineHeader)
}

type EventLineHeader struct {
	// Index original entry index in absolute file offset
	Index int64
	// Index original entry index in absolute file offset, after Layer field onwards
	IndexFromLayer int64
	// Length of line in bytes
	Length int
	// Length of line in bytes, without Layer field
	LengthFromLayer int

	// Depth The placement in the display list, including nesting
	Depth types.Depth
	// Start When the frame appears
	Start time.Duration

	ReadOrder int
}

func SortEventLineHeader(a, b EventLineHeader) int {
	//TODO: check order
	depthCmp := a.Depth.Compare(b.Depth)
	if depthCmp != 0 {
		return depthCmp
	}
	//TODO: check order
	startCmp := int(a.Start - b.Start)
	if startCmp != 0 {
		return startCmp
	}
	return a.ReadOrder - b.ReadOrder
}
