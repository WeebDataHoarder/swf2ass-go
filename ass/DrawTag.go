package ass

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	"strings"
)

type DrawTag struct {
	BaseDrawingTag
	Scale int64
}

func NewDrawTag(shape *shapes.Shape, scale int64) *DrawTag {
	return &DrawTag{
		Scale:          scale,
		BaseDrawingTag: BaseDrawingTag(*shape),
	}
}

func (t *DrawTag) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) DrawingTag {
	return &DrawTag{
		BaseDrawingTag: BaseDrawingTag(*t.AsShape().ApplyMatrixTransform(transform, applyTranslation)),
		Scale:          t.Scale,
	}
}

func (t *DrawTag) TransitionShape(line *Line, shape *shapes.Shape) PathTag {
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
	scaleMultiplier := int64(1 << (t.Scale - 1))
	precision := DefaultDrawingPrecision
	if t.Scale >= 5 {
		precision = 0
	}
	return fmt.Sprintf("\\p%d}%s{\\p0", t.Scale, strings.Join(t.GetCommands(scaleMultiplier, int64(precision)), " "))
}
