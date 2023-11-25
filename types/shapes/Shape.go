package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/records"
)

type Shape struct {
	Edges []records.Record

	IsFlat bool
}

func NewShape(edges []records.Record) *Shape {
	s := &Shape{
		IsFlat: true,
	}
	s.Edges = make([]records.Record, 0, len(edges))
	for i := range edges {
		s.AddRecord(edges[i])
	}
	return s
}

func (s *Shape) AddRecord(record records.Record) {
	if !record.IsFlat() {
		s.IsFlat = false
	}

	s.Edges = append(s.Edges, record)
}

func (s *Shape) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) *Shape {
	newShape := NewShape(nil)
	newShape.Edges = make([]records.Record, 0, len(s.Edges))
	for _, edge := range s.Edges {
		newShape.AddRecord(edge.ApplyMatrixTransform(transform, applyTranslation))
	}
	return newShape
}

func (s *Shape) Start() math.Vector2[float64] {
	if len(s.Edges) == 0 {
		return math.NewVector2[float64](0, 0)
	}
	return s.Edges[0].GetStart()
}

func (s *Shape) End() math.Vector2[float64] {
	if len(s.Edges) == 0 {
		return math.NewVector2[float64](0, 0)
	}
	return s.Edges[len(s.Edges)-1].GetEnd()
}

func (s *Shape) IsClosed() bool {
	return s.Start().Equals(s.End())
}

func (s *Shape) Reverse() *Shape {
	r := &Shape{
		Edges:  make([]records.Record, len(s.Edges)),
		IsFlat: s.IsFlat,
	}
	for i, e := range s.Edges {
		r.Edges[len(s.Edges)-1-i] = e.Reverse()
	}
	return r
}

func (s *Shape) Merge(o *Shape) *Shape {
	r := &Shape{
		Edges: make([]records.Record, 0, len(s.Edges)+len(o.Edges)),
	}
	if s.IsFlat == o.IsFlat {
		r.IsFlat = s.IsFlat
	}
	r.Edges = append(r.Edges, s.Edges...)
	r.Edges = append(r.Edges, o.Edges...)
	return r
}

// Flatten Converts all non-linear records into line segments and returns a new Shape
func (s *Shape) Flatten() *Shape {
	if s.IsFlat {
		return s
	}
	r := &Shape{
		Edges:  make([]records.Record, 0, len(s.Edges)*4),
		IsFlat: true,
	}

	for _, e := range s.Edges {
		if !e.IsFlat() {
			switch ce := e.(type) {
			case *records.QuadraticCurveRecord:
				for _, lr := range ce.ToLineRecords(1) {
					rec := lr
					r.Edges = append(r.Edges, rec)
				}
			case *records.CubicCurveRecord:
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

func IterateMorphShape(start, end *Shape) (r []records.RecordPair) {

	startEdges := start.Edges
	endEdges := end.Edges

	var prevStart, prevEnd records.Record

	for len(startEdges) > 0 && len(endEdges) > 0 {
		startEdge := startEdges[0]
		endEdge := endEdges[0]

		advanceStart := true
		advanceEnd := true

		if prevStart != nil && !prevStart.GetEnd().Equals(startEdge.GetStart()) {
			advanceStart = false
			startEdge = &records.MoveRecord{
				To:    startEdge.GetStart(),
				Start: prevStart.GetEnd(),
			}
		}

		if prevEnd != nil && !prevEnd.GetEnd().Equals(endEdge.GetStart()) {
			advanceEnd = false
			endEdge = &records.MoveRecord{
				To:    endEdge.GetStart(),
				Start: prevEnd.GetEnd(),
			}
		}

		if startEdge.SameType(endEdge) {
			r = append(r, records.RecordPair{startEdge, endEdge})
		} else {
			aLineRecord, aIsLineRecord := startEdge.(*records.LineRecord)
			aMoveRecord, aIsMoveRecord := startEdge.(*records.MoveRecord)
			aQuadraticCurveRecord, aIsQuadraticCurveRecord := startEdge.(*records.QuadraticCurveRecord)
			bLineRecord, bIsLineRecord := endEdge.(*records.LineRecord)
			bMoveRecord, bIsMoveRecord := endEdge.(*records.MoveRecord)
			bQuadraticCurveRecord, bIsQuadraticCurveRecord := endEdge.(*records.QuadraticCurveRecord)

			if aIsLineRecord && bIsQuadraticCurveRecord {
				startEdge = records.QuadraticCurveFromLineRecord(aLineRecord)
				r = append(r, records.RecordPair{startEdge, bQuadraticCurveRecord})
			} else if aIsQuadraticCurveRecord && bIsLineRecord {
				endEdge = records.QuadraticCurveFromLineRecord(bLineRecord)
				r = append(r, records.RecordPair{aQuadraticCurveRecord, endEdge})
			} else if aIsMoveRecord && !bIsMoveRecord {
				endEdge = &records.MoveRecord{
					To:    endEdge.GetStart(),
					Start: endEdge.GetStart(),
				}
				r = append(r, records.RecordPair{aMoveRecord, endEdge})
				advanceEnd = false
			} else if !aIsMoveRecord && bIsMoveRecord {
				startEdge = &records.MoveRecord{
					To:    startEdge.GetStart(),
					Start: startEdge.GetStart(),
				}
				r = append(r, records.RecordPair{startEdge, bMoveRecord})
				advanceStart = false
			} else {
				panic("incompatible")
			}
		}

		if advanceStart {
			startEdges = startEdges[1:]
		}

		if advanceEnd {
			endEdges = endEdges[1:]
		}

		prevStart = startEdge
		prevEnd = endEdge
	}

	if len(startEdges) != 0 || len(endEdges) != 0 {
		panic("incompatible result")
	}

	return r
}
