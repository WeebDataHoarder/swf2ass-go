package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/records"
	"slices"
)

type PathSegment[T ~float64 | ~int64] []VisitedPoint[T]

func NewPathSegment[T ~float64 | ~int64](start math.Vector2[T]) PathSegment[T] {
	return PathSegment[T]{
		VisitedPoint[T]{
			Pos:             start,
			IsBezierControl: false,
		},
	}
}

// Flip
// Flips the direction of the path segment.
// Flash fill paths are dual-sided, with fill style 1 indicating the positive side
// and fill style 0 indicating the negative. We have to flip fill style 0 paths
// in order to link them to fill style 1 paths.
func (s *PathSegment[T]) Flip() {
	slices.Reverse(*s)
}

func (s *PathSegment[T]) AddPoint(p VisitedPoint[T]) {
	*s = append(*s, p)
}

func (s *PathSegment[T]) Start() math.Vector2[T] {
	return (*s)[0].Pos
}

func (s *PathSegment[T]) End() math.Vector2[T] {
	return (*s)[len(*s)-1].Pos
}

func (s *PathSegment[T]) IsEmpty() bool {
	return len(*s) <= 1
}

func (s *PathSegment[T]) IsClosed() bool {
	return s.Start().Equals(s.End())
}

func (s *PathSegment[T]) Swap(o *PathSegment[T]) {
	*s, *o = *o, *s
}

func (s *PathSegment[T]) Merge(o PathSegment[T]) {
	*s = append(*s, o[1:]...)
}

func (s *PathSegment[T]) TryMerge(o *PathSegment[T], isDirected bool) bool {
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

func (s *PathSegment[T]) GetShape() (shape Shape) {
	if s.IsEmpty() {
		panic("not possible")
	}

	shape = make(Shape, 0, len(*s)-1)

	points := *s

	next := func() VisitedPoint[T] {
		point := points[0]
		points = points[1:]
		return point
	}

	lastPos := next().Pos.Float64()
	//lastPos := points[0].Pos.Float64()

	for len(points) > 0 {
		point := next()

		if !point.IsBezierControl {
			shape = append(shape, records.LineRecord{
				To:    point.Pos.Float64(),
				Start: lastPos,
			})
			lastPos = point.Pos.Float64()
		} else {
			if len(points) == 0 {
				panic("Bezier without endpoint")
			}
			end := next()

			shape = append(shape, records.QuadraticCurveRecord{
				Control: point.Pos.Float64(),
				Anchor:  end.Pos.Float64(),
				Start:   lastPos,
			})
			lastPos = end.Pos.Float64()
		}
	}

	return shape
}
