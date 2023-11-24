package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type BlurTag struct {
	Blur int64
}

func (t *BlurTag) FromStyleRecord(record shapes.StyleRecord, transform math.MatrixTransform) StyleTag {
	if lineStyleRecord, ok := record.(*shapes.LineStyleRecord); ok {
		t.Blur = int64(lineStyleRecord.Blur)
	} else if fillStyleRecord, ok := record.(*shapes.FillStyleRecord); ok {
		t.Blur = int64(fillStyleRecord.Blur)
	} else {
		t.Blur = 0
	}
	return t
}

func (t *BlurTag) TransitionStyleRecord(event Event, record shapes.StyleRecord, transform math.MatrixTransform) StyleTag {
	t2 := &BlurTag{}
	t2.FromStyleRecord(record, transform)
	return t2
}

func (t *BlurTag) Equals(tag Tag) bool {
	if o, ok := tag.(*BlurTag); ok {
		return *t == *o
	}
	return false
}

func (t *BlurTag) Encode(event time.EventTime) string {
	return fmt.Sprintf("\\blur%d", t.Blur)
}
