package records

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	math2 "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"math"
)

type CubicCurveRecord struct {
	Control1, Control2 math2.Vector2[types.Twip]
	Anchor             math2.Vector2[types.Twip]
	Start              math2.Vector2[types.Twip]
}

func (r *CubicCurveRecord) GetStart() math2.Vector2[types.Twip] {
	return r.Start
}

func (r *CubicCurveRecord) GetEnd() math2.Vector2[types.Twip] {
	return r.Anchor
}

func (r *CubicCurveRecord) Reverse() Record {
	return &CubicCurveRecord{
		Control1: r.Control2,
		Control2: r.Control1,
		Anchor:   r.Start,
		Start:    r.Anchor,
	}
}

func (r *CubicCurveRecord) ApplyMatrixTransform(transform math2.MatrixTransform, applyTranslation bool) Record {
	//TODO: see how accurate this is
	return &CubicCurveRecord{
		Control1: math2.MatrixTransformApplyToVector(transform, r.Control1, applyTranslation),
		Control2: math2.MatrixTransformApplyToVector(transform, r.Control2, applyTranslation),
		Anchor:   math2.MatrixTransformApplyToVector(transform, r.Anchor, applyTranslation),
		Start:    math2.MatrixTransformApplyToVector(transform, r.Start, applyTranslation),
	}
}

func (r *CubicCurveRecord) Equals(other Record) bool {
	if o, ok := other.(*CubicCurveRecord); ok {
		return *o == *r
	}
	return false
}

func (r *CubicCurveRecord) SameType(other Record) bool {
	_, ok := other.(*CubicCurveRecord)
	return ok
}

func (r *CubicCurveRecord) IsFlat() bool {
	return false
}

func CubicCurveFromQuadraticRecord(q *QuadraticCurveRecord) *CubicCurveRecord {
	return &CubicCurveRecord{
		Control1: q.Start.AddVector(q.Control.Multiply(2)).Divide(3),
		Control2: q.Anchor.AddVector(q.Control.Multiply(2)).Divide(3),
		Anchor:   q.Anchor,
		Start:    q.Start,
	}
}

// ToSingleQuadraticRecord Finds if Cubic curve is a perfect fit of a Quadratic curve (aka, it was upconverted)
func (r *CubicCurveRecord) ToSingleQuadraticRecord() *QuadraticCurveRecord {
	control1 := r.Control1.Multiply(3).SubVector(r.Start).Divide(2)
	control2 := r.Control2.Multiply(3).SubVector(r.Anchor).Divide(2)
	if control1.Equals(control2) {
		return &QuadraticCurveRecord{
			Control: control1,
			Anchor:  r.Anchor,
			Start:   r.Start,
		}
	}
	return nil
}

func (r *CubicCurveRecord) ToLineRecords(scale int64) []*LineRecord {
	distanceToleranceSquare := math.Pow(0.5/float64(scale), 2)
	points := CubicRecursiveBezier(nil, 0.0, 0.0, distanceToleranceSquare, r.Start.Float64(), r.Control1.Float64(), r.Control2.Float64(), r.Anchor.Float64(), 0)

	result := make([]*LineRecord, 0, len(points)+1)

	var current = r.Start

	for _, point := range points {
		tp := math2.Vector2ToType[float64, types.Twip](point)
		result = append(result, &LineRecord{
			To:    tp,
			Start: current,
		})
		current = tp
	}

	result = append(result, &LineRecord{
		To:    r.Anchor,
		Start: current,
	})

	return result
}
