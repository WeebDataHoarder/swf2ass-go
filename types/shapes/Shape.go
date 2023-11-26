package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/records"
)

type Shape []records.Record

func (s Shape) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) (newShape Shape) {
	newShape = make(Shape, 0, len(s))
	for _, edge := range s {
		newShape = append(newShape, edge.ApplyMatrixTransform(transform, applyTranslation))
	}
	return newShape
}

func (s Shape) Start() math.Vector2[float64] {
	if len(s) == 0 {
		return math.NewVector2[float64](0, 0)
	}
	return s[0].GetStart()
}

func (s Shape) End() math.Vector2[float64] {
	if len(s) == 0 {
		return math.NewVector2[float64](0, 0)
	}
	return s[len(s)-1].GetEnd()
}

func (s Shape) IsFlat() bool {
	for _, e := range s {
		if !e.IsFlat() {
			return false
		}
	}
	return true
}

func (s Shape) IsClosed() bool {
	return s.Start().Equals(s.End())
}

func (s Shape) Reverse() (r Shape) {
	r = make(Shape, len(s))
	for i, e := range s {
		r[len(s)-1-i] = e.Reverse()
	}
	return r
}

func (s Shape) Merge(o Shape) (r Shape) {
	r = make(Shape, 0, len(s)+len(o))
	r = append(r, s...)
	r = append(r, o...)
	return r
}

// Flatten Converts all non-linear records into line segments and returns a new Shape
func (s Shape) Flatten() (r Shape) {
	if s.IsFlat() {
		return s
	}
	r = make(Shape, 0, len(s)*4)

	for _, e := range s {
		r = append(r, records.FlattenRecord(e, 1)...)
	}

	return r
}

func (s Shape) Equals(o Shape) bool {
	if len(s) != len(o) {
		return false
	}

	for i := range s {
		if !s[i].Equals(o[i]) {
			return false
		}
	}
	return true
}

func IterateMorphShape(start, end Shape) (r []records.RecordPair) {

	var prevStart, prevEnd records.Record

	for len(start) > 0 && len(end) > 0 {
		startEdge := start[0]
		endEdge := end[0]

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
			start = start[1:]
		}

		if advanceEnd {
			end = end[1:]
		}

		prevStart = startEdge
		prevEnd = endEdge
	}

	if len(start) != 0 || len(end) != 0 {
		panic("incompatible result")
	}

	return r
}
