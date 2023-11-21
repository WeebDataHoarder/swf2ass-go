package records

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	math2 "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"math"
)

type QuadraticCurveRecord struct {
	Control math2.Vector2[types.Twip]
	Anchor  math2.Vector2[types.Twip]
	Start   math2.Vector2[types.Twip]
}

func (r *QuadraticCurveRecord) GetStart() math2.Vector2[types.Twip] {
	return r.Start
}

func (r *QuadraticCurveRecord) GetEnd() math2.Vector2[types.Twip] {
	return r.Anchor
}

func (r *QuadraticCurveRecord) Reverse() Record {
	return &QuadraticCurveRecord{
		Control: r.Control,
		Anchor:  r.Start,
		Start:   r.Anchor,
	}
}

func (r *QuadraticCurveRecord) ApplyMatrixTransform(transform math2.MatrixTransform, applyTranslation bool) Record {
	//TODO: see how accurate this is
	return &QuadraticCurveRecord{
		Control: math2.Vector2ToType[float64, types.Twip](transform.ApplyToVector(r.Control.Float64().Divide(types.TwipFactor), applyTranslation).Multiply(types.TwipFactor)),
		Anchor:  math2.Vector2ToType[float64, types.Twip](transform.ApplyToVector(r.Anchor.Float64().Divide(types.TwipFactor), applyTranslation).Multiply(types.TwipFactor)),
		Start:   math2.Vector2ToType[float64, types.Twip](transform.ApplyToVector(r.Start.Float64().Divide(types.TwipFactor), applyTranslation).Multiply(types.TwipFactor)),
	}
}

func (r *QuadraticCurveRecord) Equals(other Record) bool {
	if o, ok := other.(*QuadraticCurveRecord); ok {
		return *o == *r
	}
	return false
}

func (r *QuadraticCurveRecord) SameType(other Record) bool {
	_, ok := other.(*QuadraticCurveRecord)
	return ok
}

func (r *QuadraticCurveRecord) IsFlat() bool {
	return false
}

func QuadraticCurveFromLineRecord(l *LineRecord) *QuadraticCurveRecord {
	delta := l.To.SubVector(l.Start)
	return &QuadraticCurveRecord{
		Control: l.Start.AddVector(delta.Divide(2)),
		Anchor:  l.Start.AddVector(delta),
		Start:   l.Start,
	}
}

func (r *QuadraticCurveRecord) ToLineRecords(scale int64) []*LineRecord {
	distanceToleranceSquare := math.Pow(0.5/float64(scale), 2)
	points := QuadraticRecursiveBezier(nil, 0.0, distanceToleranceSquare, r.Start.Float64().Divide(types.TwipFactor), r.Control.Float64().Divide(types.TwipFactor), r.Anchor.Float64().Divide(types.TwipFactor), 0)

	result := make([]*LineRecord, 0, len(points)+1)

	var current = r.Start

	for _, point := range points {
		tp := math2.Vector2ToType[float64, types.Twip](point.Multiply(types.TwipFactor))
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
