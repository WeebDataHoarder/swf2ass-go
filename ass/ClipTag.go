package ass

import (
	"fmt"
	swftypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"strings"
)

type ClipTag struct {
	BaseDrawingTag
	Scale  int64
	IsNull bool
}

func NewClipTag(clip *types.ClipPath, scale int64) *ClipTag {
	if clip == nil {
		return &ClipTag{
			IsNull: true,
			Scale:  scale,
		}
	} else {
		shape := clip.GetShape()
		if len(shape.Edges) == 0 { //full clip
			shape = &types.Shape{
				Edges: []types.Record{
					&types.LineRecord{
						//TODO: ??? why swftypes.TwipFactor here???
						To:    types.NewVector2[swftypes.Twip](0, swftypes.TwipFactor),
						Start: types.NewVector2[swftypes.Twip](0, 0),
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

func (t *ClipTag) ApplyMatrixTransform(transform types.MatrixTransform, applyTranslation bool) DrawingTag {
	return &ClipTag{
		BaseDrawingTag: BaseDrawingTag(*transform.ApplyToShape(t.AsShape(), applyTranslation)),
	}
}

func (t *ClipTag) TransitionClipPath(line *Line, clip *types.ClipPath) ClipPathTag {
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

func (t *ClipTag) Encode(event EventTime) string {
	scaleMultiplier := 1 << t.Scale
	if t.IsNull {
		return ""
	}
	precision := DefaultDrawingPrecision
	if t.Scale >= 5 {
		precision = 0
	}
	return fmt.Sprintf("\\clip(%d,%s)", t.Scale, strings.Join(t.GetCommands(scaleMultiplier, int64(precision)), " "))
}
