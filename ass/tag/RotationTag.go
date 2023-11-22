package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	math2 "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"math"
)

type RotationTag struct {
	RotationX, RotationY, RotationZ float64
}

func (t *RotationTag) TransitionMatrixTransform(event Event, transform math2.MatrixTransform) PositioningTag {
	panic("not implemented")
}

func (t *RotationTag) Encode(event time.EventTime) string {
	//TODO: precision?
	return fmt.Sprintf("\\frx%.2F\\fry%.2F\\frz%.2F", t.RotationX, t.RotationY, t.RotationZ)
}

func (t *RotationTag) Equals(tag Tag) bool {
	if o, ok := tag.(*RotationTag); ok {
		return math.Abs(t.RotationX-o.RotationX) <= math.SmallestNonzeroFloat64 && math.Abs(t.RotationY-o.RotationY) <= math.SmallestNonzeroFloat64 && math.Abs(t.RotationZ-o.RotationZ) <= math.SmallestNonzeroFloat64
	}
	return false
}

func (t *RotationTag) FromMatrixTransform(transform math2.MatrixTransform) PositioningTag {
	//maybe qr decomposition?
	panic("not implemented")
}
