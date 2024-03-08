package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/settings"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	"github.com/ctessum/polyclip-go"
)

type ClipTag struct {
	BaseDrawingTag
	Scale  int
	IsNull bool
}

func NewClipTag(clip types.Option[shapes.Shape], scale int) *ClipTag {
	if c, ok := clip.Some(); ok && len(c) > 0 {
		return &ClipTag{
			Scale:          scale,
			BaseDrawingTag: BaseDrawingTag(c),
			IsNull:         len(c) == 0,
		}
	} else {
		return &ClipTag{
			IsNull: true,
			Scale:  scale,
		}
	}
}

func (t *ClipTag) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) DrawingTag {
	return &ClipTag{
		BaseDrawingTag: BaseDrawingTag(t.AsShape().ApplyMatrixTransform(transform, applyTranslation)),
		Scale:          t.Scale,
	}
}

func (t *ClipTag) TransitionClipPath(event Event, clip *shapes.ClipPath) ClipPathTag {
	if clip == nil {
		if t.IsNull {
			return t
		} else {
			return nil
		}
	}
	if t.AsShape().Equals(clip.GetShape()) {
		return t
	} else {
		return nil
	}
}

func (t *ClipTag) Equals(tag Tag) bool {
	if o, ok := tag.(*ClipTag); ok {
		return t.IsNull == o.IsNull && t.Scale == o.Scale && t.AsShape().Equals(o.AsShape())
	}
	return false
}

func (t *ClipTag) Encode(event time.EventTime) string {
	if t.IsNull {
		return ""
	}

	shape := t.AsShape()
	bb := shape.BoundingBox()
	//uses pixel coords
	if bb.TopLeft.Int64().Float64().Equals(bb.TopLeft) && bb.BottomRight.Int64().Float64().Equals(bb.BottomRight) {
		diffPol := shapes.NewPolygonFromShape(bb.Draw()).Construct(polyclip.DIFFERENCE, shapes.NewPolygonFromShape(shape))
		if len(diffPol) == 0 { //it's the same!
			//we can use square clip!
			return fmt.Sprintf("\\clip(%d,%d,%d,%d)", bb.TopLeft.Int64().X, bb.TopLeft.Int64().Y, bb.BottomRight.Int64().X, bb.BottomRight.Int64().Y)
		}
	}

	scaleMultiplier := 1 << (t.Scale - 1)
	precision := settings.GlobalSettings.ASSDrawingPrecision
	if t.Scale >= 5 {
		precision = 0
	}
	return fmt.Sprintf("\\clip(%d,%s)", t.Scale, t.GetCommands(scaleMultiplier, precision))
}
