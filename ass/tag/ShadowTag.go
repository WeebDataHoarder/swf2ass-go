package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/line"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	swftypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type ShadowTag struct {
	Depth swftypes.Twip
}

func (t *ShadowTag) FromStyleRecord(record shapes.StyleRecord) StyleTag {
	//TODO?
	t.Depth = 0
	return t
}

func (t *ShadowTag) TransitionStyleRecord(line *line.Line, record shapes.StyleRecord) StyleTag {
	t2 := &ShadowTag{}
	t2.FromStyleRecord(record)
	return t2
}

func (t *ShadowTag) Equals(tag Tag) bool {
	if o, ok := tag.(*ShadowTag); ok {
		return *t == *o
	}
	return false
}

func (t *ShadowTag) Encode(event time.EventTime) string {
	return fmt.Sprintf("\\shad%.02F", t.Depth.Float64())
}
