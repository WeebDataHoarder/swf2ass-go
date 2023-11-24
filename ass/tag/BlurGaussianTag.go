package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type BlurGaussianTag struct {
	Blur float64
}

func (t *BlurGaussianTag) FromStyleRecord(record shapes.StyleRecord, transform math.MatrixTransform) StyleTag {
	if lineStyleRecord, ok := record.(*shapes.LineStyleRecord); ok {
		t.Blur = lineStyleRecord.Blur
	} else if fillStyleRecord, ok := record.(*shapes.FillStyleRecord); ok {
		t.Blur = fillStyleRecord.Blur
	} else {
		t.Blur = 0
	}
	return t
}

func (t *BlurGaussianTag) TransitionStyleRecord(event Event, record shapes.StyleRecord, transform math.MatrixTransform) StyleTag {
	t2 := &BlurGaussianTag{}
	t2.FromStyleRecord(record, transform)
	return t2
}

func (t *BlurGaussianTag) Equals(tag Tag) bool {
	if o, ok := tag.(*BlurGaussianTag); ok {
		return *t == *o
	}
	return false
}

func (t *BlurGaussianTag) Encode(event time.EventTime) string {
	if t.Blur == 0 {
		return "\\blur0"
	}
	return fmt.Sprintf("\\blur%.02F", t.Blur)
}
