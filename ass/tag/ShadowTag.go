package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type ShadowTag struct {
	Depth float64
}

func (t *ShadowTag) FromStyleRecord(record shapes.StyleRecord, transform math.MatrixTransform) StyleTag {
	//TODO?
	t.Depth = 0
	return t
}

func (t *ShadowTag) TransitionStyleRecord(event Event, record shapes.StyleRecord, transform math.MatrixTransform) StyleTag {
	t2 := &ShadowTag{}
	t2.FromStyleRecord(record, transform)
	return t2
}

func (t *ShadowTag) Equals(tag Tag) bool {
	if o, ok := tag.(*ShadowTag); ok {
		return *t == *o
	}
	return false
}

func (t *ShadowTag) Encode(event time.EventTime) string {
	if t.Depth == 0 {
		return "\\shad0"
	}

	return fmt.Sprintf("\\shad%.02F", t.Depth)
}
