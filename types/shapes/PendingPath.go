package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"slices"
)

type PendingPathMap map[int]*PendingPath[types.Twip]

func (m PendingPathMap) MergePath(p *ActivePath, directed bool) {
	if _, ok := m[p.StyleId]; !ok {
		m[p.StyleId] = &PendingPath[types.Twip]{}
	}
	m[p.StyleId].MergePath(&p.Segment, directed)
}

type PendingPath[T ~float64 | ~int64] []*PathSegment[T]

func (p *PendingPath[T]) MergePath(newSegment *PathSegment[T], directed bool) {
	if !newSegment.IsEmpty() {
		var merged *PathSegment[T]

		for i, segment := range *p {
			if segment.TryMerge(newSegment, directed) {
				*p = slices.Delete(*p, i, i+1)
				merged = segment
				break
			}
		}

		if merged != nil {
			p.MergePath(merged, directed)
		} else {
			*p = append(*p, newSegment)
		}
	}
}

func (p *PendingPath[T]) GetShape() *Shape {
	shape := &Shape{}
	for _, segment := range *p {
		shape = shape.Merge(segment.GetShape())
	}
	return shape
}
