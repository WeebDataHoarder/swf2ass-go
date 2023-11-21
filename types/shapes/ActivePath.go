package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type ActivePath struct {
	Segment PathSegment
	StyleId int
}

func NewActivePath(styleId int, start math.Vector2[types.Twip]) *ActivePath {
	return &ActivePath{
		Segment: NewPathSegment(start),
		StyleId: styleId,
	}
}

func (p *ActivePath) AddPoint(point VisitedPoint) {
	p.Segment.AddPoint(point)
}

func (p *ActivePath) Flip() {
	p.Segment.Flip()
}
