package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/records"
	"slices"
)

type PathSegment []VisitedPoint

func NewPathSegment(start math.Vector2[types.Twip]) PathSegment {
	return PathSegment{
		VisitedPoint{
			Pos:             start,
			IsBezierControl: false,
		},
	}
}

func (s *PathSegment) Flip() {
	slices.Reverse(*s)
}

func (s *PathSegment) AddPoint(p VisitedPoint) {
	*s = append(*s, p)
}

func (s *PathSegment) Start() math.Vector2[types.Twip] {
	return (*s)[0].Pos
}

func (s *PathSegment) End() math.Vector2[types.Twip] {
	return (*s)[len(*s)-1].Pos
}

func (s *PathSegment) IsEmpty() bool {
	return len(*s) <= 1
}

func (s *PathSegment) IsClosed() bool {
	return s.Start().Equals(s.End())
}

func (s *PathSegment) Swap(o *PathSegment) {
	*s, *o = *o, *s
}

func (s *PathSegment) Merge(o PathSegment) {
	*s = append(*s, o[1:]...)
}

func (s *PathSegment) TryMerge(o *PathSegment, isDirected bool) bool {
	if o.End().Equals(s.Start()) {
		s.Swap(o)
		s.Merge(*o)
		return true
	} else if s.End().Equals(o.Start()) {
		s.Merge(*o)
		return true
	} else if !isDirected && s.End().Equals(o.End()) {
		o.Flip()
		s.Merge(*o)
		return true
	} else if !isDirected && s.Start().Equals(o.Start()) {
		o.Flip()
		s.Swap(o)
		s.Merge(*o)
		return true
	}

	return false
}

func (s *PathSegment) GetShape() *Shape {
	if s.IsEmpty() {
		panic("not possible")
	}

	shape := &Shape{
		Edges: make([]records.Record, 0, len(*s)),
	}

	pos := s.Start()

	points := *s
	for len(points) > 0 {
		point := points[0]
		points = points[1:]
		if !point.IsBezierControl {
			shape.AddRecord(&records.LineRecord{
				To:    point.Pos,
				Start: pos,
			})
			pos = point.Pos
		} else {
			if len(points) == 0 {
				panic("bezier without endpoint")
			}
			end := points[0]
			points = points[1:]

			shape.AddRecord(&records.QuadraticCurveRecord{
				Control: point.Pos,
				Anchor:  end.Pos,
				Start:   pos,
			})
			pos = end.Pos
		}
	}

	return shape
}
