package shapes

import "slices"

type PendingPathMap map[int]*PendingPath

func (m PendingPathMap) MergePath(p *ActivePath, directed bool) {
	if _, ok := m[p.StyleId]; !ok {
		m[p.StyleId] = &PendingPath{}
	}
	m[p.StyleId].MergePath(&p.Segment, directed)
}

type PendingPath []*PathSegment

func (p *PendingPath) MergePath(newSegment *PathSegment, directed bool) {
	if !newSegment.IsEmpty() {
		var merged *PathSegment

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

func (p *PendingPath) GetShape() *Shape {
	shape := &Shape{}
	for _, segment := range *p {
		shape = shape.Merge(segment.GetShape())
	}
	return shape
}
