package types

import "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"

type Shape struct {
	Edges []Record

	IsFlat bool
}

func NewShape(edges []Record) *Shape {
	s := &Shape{
		IsFlat: true,
	}
	s.Edges = make([]Record, 0, len(edges))
	for i := range edges {
		s.AddRecord(edges[i])
	}
	return s
}

func (s *Shape) AddRecord(record Record) {
	if !record.IsFlat() {
		s.IsFlat = false
	}

	s.Edges = append(s.Edges, record)
}

func (s *Shape) Start() Vector2[types.Twip] {
	if len(s.Edges) == 0 {
		return NewVector2[types.Twip](0, 0)
	}
	return s.Edges[0].GetStart()
}

func (s *Shape) End() Vector2[types.Twip] {
	if len(s.Edges) == 0 {
		return NewVector2[types.Twip](0, 0)
	}
	return s.Edges[len(s.Edges)-1].GetEnd()
}

func (s *Shape) IsClosed() bool {
	return s.Start().Equals(s.End())
}

func (s *Shape) Merge(o *Shape) *Shape {
	r := &Shape{
		Edges: make([]Record, 0, len(s.Edges)+len(o.Edges)),
	}
	if s.IsFlat == o.IsFlat {
		r.IsFlat = s.IsFlat
	}
	r.Edges = append(r.Edges, s.Edges...)
	r.Edges = append(r.Edges, o.Edges...)
	return r
}

func (s *Shape) Flatten() *Shape {
	if s.IsFlat {
		return s
	}
	r := &Shape{
		Edges:  make([]Record, 0, len(s.Edges)*4),
		IsFlat: true,
	}

	for _, e := range s.Edges {
		if !e.IsFlat() {
			switch ce := e.(type) {
			case *QuadraticCurveRecord:
				for _, lr := range ce.ToLineRecords(1) {
					rec := lr
					r.Edges = append(r.Edges, rec)
				}
			case *CubicCurveRecord:
				for _, lr := range ce.ToLineRecords(1) {
					rec := lr
					r.Edges = append(r.Edges, rec)
				}
			default:
				panic("not implemented")
			}
		} else {
			r.Edges = append(r.Edges, e)
		}
	}

	return r
}

func (s *Shape) Equals(o *Shape) bool {
	if len(s.Edges) != len(o.Edges) && s.IsFlat == o.IsFlat /* todo: check this last condition */ {
		return false
	}

	for i := range s.Edges {
		if !s.Edges[i].Equals(o.Edges[i]) {
			return false
		}
	}
	return true
}
