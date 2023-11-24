package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"strconv"
)

type ScaleTag struct {
	Scale math.Vector2[float64]
}

func (t *ScaleTag) TransitionMatrixTransform(event Event, transform math.MatrixTransform) PositioningTag {
	panic("not implemented")
}

func (t *ScaleTag) Encode(event time.EventTime) string {
	//TODO: precision?
	return fmt.Sprintf("\\fscx%s\\fscy%s", strconv.FormatFloat(t.Scale.X, 'f', -1, 64), strconv.FormatFloat(t.Scale.Y, 'f', -1, 64))
}

func (t *ScaleTag) Equals(tag Tag) bool {
	if o, ok := tag.(*ScaleTag); ok {
		return t.Scale.Equals(o.Scale)
	}
	return false
}

func (t *ScaleTag) FromMatrixTransform(transform math.MatrixTransform) PositioningTag {
	//maybe qr decomposition?
	panic("not implemented")
}
