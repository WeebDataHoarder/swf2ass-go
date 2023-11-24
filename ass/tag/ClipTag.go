package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/settings"
	swftypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/records"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type ClipTag struct {
	BaseDrawingTag
	Scale  int
	IsNull bool
}

func NewClipTag(clip *shapes.ClipPath, scale int) *ClipTag {
	if clip == nil {
		return &ClipTag{
			IsNull: true,
			Scale:  scale,
		}
	} else {
		shape := clip.GetShape()
		if len(shape.Edges) == 0 { //full clip
			shape = &shapes.Shape{
				Edges: []records.Record{
					&records.LineRecord{
						//TODO: ??? why TwipFactor here???
						To:    math.NewVector2[float64](0, swftypes.Twip(swftypes.TwipFactor).Float64()),
						Start: math.NewVector2[float64](0, 0),
					},
				},
				IsFlat: true,
			}
		}
		return &ClipTag{
			Scale:          scale,
			BaseDrawingTag: BaseDrawingTag(*shape),
		}
	}
}

func (t *ClipTag) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) DrawingTag {
	return &ClipTag{
		BaseDrawingTag: BaseDrawingTag(*t.AsShape().ApplyMatrixTransform(transform, applyTranslation)),
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
	scaleMultiplier := 1 << (t.Scale - 1)
	precision := settings.GlobalSettings.ASSDrawingPrecision
	if t.Scale >= 5 {
		precision = 0
	}
	return fmt.Sprintf("\\clip(%d,%s)", t.Scale, t.GetCommands(scaleMultiplier, precision))
}
