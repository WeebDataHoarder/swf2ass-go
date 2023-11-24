package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"strconv"
)

type ShearingTag struct {
	Shear math.Vector2[float64]
}

func (t *ShearingTag) TransitionMatrixTransform(event Event, transform math.MatrixTransform) PositioningTag {
	panic("not implemented")
}

func (t *ShearingTag) Encode(event time.EventTime) string {
	//TODO: precision?
	return fmt.Sprintf("\\fax%s\\fay%s", strconv.FormatFloat(t.Shear.X, 'f', -1, 64), strconv.FormatFloat(t.Shear.Y, 'f', -1, 64))
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
