package records

import (
	"fmt"
	math2 "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"math"
)

type CubicCurveRecord struct {
	Control1, Control2 math2.Vector2[float64]
	Anchor             math2.Vector2[float64]
	Start              math2.Vector2[float64]
}

func (r CubicCurveRecord) GetStart() math2.Vector2[float64] {
	return r.Start
}

func (r CubicCurveRecord) GetEnd() math2.Vector2[float64] {
	return r.Anchor
}

func (r CubicCurveRecord) Reverse() Record {
	return CubicCurveRecord{
		Control1: r.Control2,
		Control2: r.Control1,
		Anchor:   r.Start,
		Start:    r.Anchor,
	}
}

func (r CubicCurveRecord) ApplyMatrixTransform(transform math2.MatrixTransform, applyTranslation bool) Record {
	//TODO: see how accurate this is
	return CubicCurveRecord{
		Control1: transform.ApplyToVector(r.Control1, applyTranslation),
		Control2: transform.ApplyToVector(r.Control2, applyTranslation),
		Anchor:   transform.ApplyToVector(r.Anchor, applyTranslation),
		Start:    transform.ApplyToVector(r.Start, applyTranslation),
	}
}

func (r CubicCurveRecord) Equals(other Record) bool {
	if o, ok := other.(CubicCurveRecord); ok {
		return o == r
	}
	return false
}

func (r CubicCurveRecord) SameType(other Record) bool {
	_, ok := other.(CubicCurveRecord)
	return ok
}

func (r CubicCurveRecord) IsFlat() bool {
	return false
}

func (r CubicCurveRecord) String() string {
	return fmt.Sprintf("c %s %s %s", r.Control1, r.Control2, r.Anchor)
}

func (r CubicCurveRecord) BoundingBox() (topLeft, bottomRight math2.Vector2[float64]) {
	return r.Start.Min(r.Control1).Min(r.Control2).Min(r.Anchor), r.Start.Max(r.Control1).Max(r.Control2).Max(r.Anchor)
}

func CubicCurveFromQuadraticRecord(q QuadraticCurveRecord) CubicCurveRecord {
	return CubicCurveRecord{
		Control1: q.Start.AddVector(q.Control.Multiply(2)).Divide(3),
		Control2: q.Anchor.AddVector(q.Control.Multiply(2)).Divide(3),
		Anchor:   q.Anchor,
		Start:    q.Start,
	}
}

// ToSingleQuadraticRecord Finds if Cubic curve is a perfect fit of a Quadratic curve (aka, it was upconverted)
func (r CubicCurveRecord) ToSingleQuadraticRecord() (QuadraticCurveRecord, bool) {
	control1 := r.Control1.Multiply(3).SubVector(r.Start).Divide(2)
	control2 := r.Control2.Multiply(3).SubVector(r.Anchor).Divide(2)
	if control1.Equals(control2) {
		return QuadraticCurveRecord{
			Control: control1,
			Anchor:  r.Anchor,
			Start:   r.Start,
		}, true
	}
	return QuadraticCurveRecord{}, false
}

func (r CubicCurveRecord) ToLineRecords(scale int64) []Record {
	distanceToleranceSquare := math.Pow(0.5/float64(scale), 2)
	points := CubicRecursiveBezier(nil, 0.0, BezierCurveAngleTolerance, distanceToleranceSquare, r.Start, r.Control1, r.Control2, r.Anchor, 0)

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
