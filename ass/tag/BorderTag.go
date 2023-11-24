package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	"strconv"
)

type BorderTag struct {
	Size math.Vector2[float64]
}

func (t *BorderTag) FromStyleRecord(record shapes.StyleRecord, transform math.MatrixTransform) StyleTag {
	if lineStyleRecord, ok := record.(*shapes.LineStyleRecord); ok {
		w := lineStyleRecord.StrokeWidth(transform)
		t.Size = math.NewVector2(w, w)
	} else if fillStyleRecord, ok := record.(*shapes.FillStyleRecord); ok && fillStyleRecord.Border != nil {
		w := fillStyleRecord.Border.StrokeWidth(transform)
		t.Size = math.NewVector2(w, w)
	} else {
		t.Size = math.NewVector2[float64](0, 0)
	}
	return t
}

func (t *BorderTag) TransitionStyleRecord(event Event, record shapes.StyleRecord, transform math.MatrixTransform) StyleTag {
	t2 := &BorderTag{}
	t2.FromStyleRecord(record, transform)
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
		return fmt.Sprintf("\\bord%s", strconv.FormatFloat(t.Size.X, 'f', -1, 64))
	} else {
		return fmt.Sprintf("\\xbord%s\\ybord%s", strconv.FormatFloat(t.Size.X, 'f', -1, 64), strconv.FormatFloat(t.Size.Y, 'f', -1, 64))
	}
}
