package records

import (
	math2 "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"math"
)

type QuadraticCurveRecord struct {
	Control math2.Vector2[float64]
	Anchor  math2.Vector2[float64]
	Start   math2.Vector2[float64]
}

func (r QuadraticCurveRecord) GetStart() math2.Vector2[float64] {
	return r.Start
}

func (r QuadraticCurveRecord) GetEnd() math2.Vector2[float64] {
	return r.Anchor
}

func (r QuadraticCurveRecord) Reverse() Record {
	return QuadraticCurveRecord{
		Control: r.Control,
		Anchor:  r.Start,
		Start:   r.Anchor,
	}
}

func (r QuadraticCurveRecord) ApplyMatrixTransform(transform math2.MatrixTransform, applyTranslation bool) Record {
	//TODO: see how accurate this is
	return QuadraticCurveRecord{
		Control: math2.MatrixTransformApplyToVector(transform, r.Control, applyTranslation),
		Anchor:  math2.MatrixTransformApplyToVector(transform, r.Anchor, applyTranslation),
		Start:   math2.MatrixTransformApplyToVector(transform, r.Start, applyTranslation),
	}
}

func (r QuadraticCurveRecord) Equals(other Record) bool {
	if o, ok := other.(QuadraticCurveRecord); ok {
		return o == r
	}
	return false
}

func (r QuadraticCurveRecord) SameType(other Record) bool {
	_, ok := other.(QuadraticCurveRecord)
	return ok
}

func (r QuadraticCurveRecord) IsFlat() bool {
	return false
}

func QuadraticCurveFromLineRecord(l LineRecord) QuadraticCurveRecord {
	delta := l.To.SubVector(l.Start)
	return QuadraticCurveRecord{
		Control: l.Start.AddVector(delta.Divide(2)),
		Anchor:  l.Start.AddVector(delta),
		Start:   l.Start,
	}
}

func (r QuadraticCurveRecord) ToLineRecords(scale int64) []Record {
	distanceToleranceSquare := math.Pow(0.5/float64(scale), 2)
	points := QuadraticRecursiveBezier(nil, BezierCurveAngleTolerance, distanceToleranceSquare, r.Start, r.Control, r.Anchor, 0)

	result := make([]Record, 0, len(points)+1)

	var current = r.Start

	for _, point := range points {
		//Remove dupe segments
		if point.Equals(current) {
			continue
		}
		result = append(result, LineRecord{
			To:    point,
			Start: current,
		})
		current = point
	}

	result = append(result, LineRecord{
		To:    r.Anchor,
		Start: current,
	})

	return result
}
