package ass

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
)

type ScaleTag struct {
	Scale types.Vector2[float64]
}

func (t *ScaleTag) TransitionMatrixTransform(line *Line, transform types.MatrixTransform) PositioningTag {
	panic("not implemented")
}

func (t *ScaleTag) Encode(event EventTime) string {
	//TODO: precision?
	return fmt.Sprintf("\\fscx%.5F\\fscy%.5F", t.Scale.X, t.Scale.Y)
}

func (t *ScaleTag) Equals(tag Tag) bool {
	if o, ok := tag.(*ScaleTag); ok {
		return t.Scale.Equals(o.Scale)
	}
	return false
}

func (t *ScaleTag) FromMatrixTransform(transform types.MatrixTransform) PositioningTag {
	//maybe qr decomposition?
	panic("not implemented")
}
