package ass

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type ShearingTag struct {
	Shear math.Vector2[float64]
}

func (t *ShearingTag) TransitionMatrixTransform(line *Line, transform math.MatrixTransform) PositioningTag {
	panic("not implemented")
}

func (t *ShearingTag) Encode(event EventTime) string {
	//TODO: precision?
	return fmt.Sprintf("\\fax%.5F\\fay%.5F", t.Shear.X, t.Shear.Y)
}

func (t *ShearingTag) Equals(tag Tag) bool {
	if o, ok := tag.(*ShearingTag); ok {
		return t.Shear.Equals(o.Shear)
	}
	return false
}

func (t *ShearingTag) FromMatrixTransform(transform math.MatrixTransform) PositioningTag {
	//maybe qr decomposition?
	panic("not implemented")
}
