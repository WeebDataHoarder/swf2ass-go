package ass

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"strings"
)

type DrawTag struct {
	BaseDrawingTag
	Scale int64
}

func NewDrawTag(shape *types.Shape, scale int64) *DrawTag {
	return &DrawTag{
		Scale:          scale,
		BaseDrawingTag: BaseDrawingTag(*shape),
	}
}

func (t *DrawTag) ApplyMatrixTransform(transform types.MatrixTransform, applyTranslation bool) DrawingTag {
	return &DrawTag{
		BaseDrawingTag: BaseDrawingTag(*transform.ApplyToShape(t.AsShape(), applyTranslation)),
	}
}

func (t *DrawTag) TransitionShape(line *Line, shape *types.Shape) PathTag {
	if t.AsShape().Equals(shape) {
		return t
	}
	return nil
}

func (t *DrawTag) Equals(tag Tag) bool {
	if o, ok := tag.(*DrawTag); ok {
		return t.AsShape().Equals(o.AsShape())
	}
	return false
}

func (t *DrawTag) Encode(event EventTime) string {
	scaleMultiplier := 1 << t.Scale
	precision := DefaultDrawingPrecision
	if t.Scale >= 5 {
		precision = 0
	}
	return fmt.Sprintf("\\p%d}%s{\\p0", t.Scale, strings.Join(t.GetCommands(scaleMultiplier, int64(precision)), " "))
}
