package records

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"reflect"
	"slices"
)

type CubicSplineCurveRecord struct {
	Control []math.Vector2[float64]
	Anchor  math.Vector2[float64]
	Start   math.Vector2[float64]
}

func (r CubicSplineCurveRecord) GetStart() math.Vector2[float64] {
	return r.Start
}

func (r CubicSplineCurveRecord) GetEnd() math.Vector2[float64] {
	return r.Anchor
}

func (r CubicSplineCurveRecord) Reverse() Record {
	controls := slices.Clone(r.Control)
	slices.Reverse(controls)
	return CubicSplineCurveRecord{
		Control: controls,
		Anchor:  r.Start,
		Start:   r.Anchor,
	}
}

func (r CubicSplineCurveRecord) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) Record {
	//TODO: see how accurate this is
	controls := make([]math.Vector2[float64], 0, len(r.Control))
	for _, c := range r.Control {
		controls = append(controls, math.MatrixTransformApplyToVector(transform, c, applyTranslation))
	}
	return CubicSplineCurveRecord{
		Control: controls,
		Anchor:  math.MatrixTransformApplyToVector(transform, r.Anchor, applyTranslation),
		Start:   math.MatrixTransformApplyToVector(transform, r.Start, applyTranslation),
	}
}

func (r CubicSplineCurveRecord) Equals(other Record) bool {
	if o, ok := other.(CubicSplineCurveRecord); ok {
		return reflect.DeepEqual(r.Control, o.Control) && r.Start == o.Start && r.Anchor == o.Anchor
	}
	return false
}

func (r CubicSplineCurveRecord) SameType(other Record) bool {
	_, ok := other.(CubicSplineCurveRecord)
	return ok
}

func (r CubicSplineCurveRecord) IsFlat() bool {
	return false
}
