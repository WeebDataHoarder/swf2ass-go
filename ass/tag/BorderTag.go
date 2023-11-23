package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type BorderTag struct {
	Size math.Vector2[float64]
}

func (t *BorderTag) FromStyleRecord(record shapes.StyleRecord) StyleTag {
	if lineStyleRecord, ok := record.(*shapes.LineStyleRecord); ok {
		t.Size = math.NewVector2(lineStyleRecord.Width, lineStyleRecord.Width)
	} else if fillStyleRecord, ok := record.(*shapes.FillStyleRecord); ok && fillStyleRecord.Border != nil {
		t.Size = math.NewVector2(fillStyleRecord.Border.Width, fillStyleRecord.Border.Width)
	} else {
		t.Size = math.NewVector2[float64](0, 0)
	}
	return t
}

func (t *BorderTag) TransitionStyleRecord(event Event, record shapes.StyleRecord) StyleTag {
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

func (t *BorderTag) Encode(event time.EventTime) string {
	if t.Size.X == t.Size.Y {
		if t.Size.X == 0 {
			return "\\bord0"
		}
		return fmt.Sprintf("\\bord%.02F", t.Size.X)
	} else {
		return fmt.Sprintf("\\xbord%.02F\\ybord%.02F", t.Size.X, t.Size.Y)
	}
}
