package ass

import (
	"fmt"
	swftypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type BorderTag struct {
	Size math.Vector2[swftypes.Twip]
}

func (t *BorderTag) FromStyleRecord(record types.StyleRecord) StyleTag {
	if lineStyleRecord, ok := record.(*types.LineStyleRecord); ok {
		t.Size = math.NewVector2[swftypes.Twip](lineStyleRecord.Width, lineStyleRecord.Width)
	} else if fillStyleRecord, ok := record.(*types.FillStyleRecord); ok && fillStyleRecord.Border != nil {
		t.Size = math.NewVector2[swftypes.Twip](fillStyleRecord.Border.Width, fillStyleRecord.Border.Width)
	} else {
		t.Size = math.NewVector2[swftypes.Twip](0, 0)
	}
	return t
}

func (t *BorderTag) TransitionStyleRecord(line *Line, record types.StyleRecord) StyleTag {
	t2 := &BorderTag{}
	t2.FromStyleRecord(record)
	return t2
}

func (t *BorderTag) Equals(tag Tag) bool {
	if o, ok := tag.(*BorderTag); ok {
		return t.Size.Equals(o.Size)
	}
	return false
}

func (t *BorderTag) Encode(event EventTime) string {
	if t.Size.X == t.Size.Y {
		return fmt.Sprintf("\\bord%.02F", t.Size.X.Float64())
	} else {
		return fmt.Sprintf("\\xbord%.02F\\ybord%.02F", t.Size.X.Float64(), t.Size.Y.Float64())
	}
}
