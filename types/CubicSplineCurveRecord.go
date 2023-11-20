package types

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"reflect"
	"slices"
)

type CubicSplineCurveRecord struct {
	Control []Vector2[types.Twip]
	Anchor  Vector2[types.Twip]
	Start   Vector2[types.Twip]
}

func (r *CubicSplineCurveRecord) GetStart() Vector2[types.Twip] {
	return r.Start
}

func (r *CubicSplineCurveRecord) GetEnd() Vector2[types.Twip] {
	return r.Anchor
}

func (r *CubicSplineCurveRecord) Reverse() Record {
	controls := slices.Clone(r.Control)
	slices.Reverse(controls)
	return &CubicSplineCurveRecord{
		Control: controls,
		Anchor:  r.Start,
		Start:   r.Anchor,
	}
}

func (r *CubicSplineCurveRecord) ApplyMatrixTransform(transform MatrixTransform, applyTranslation bool) Record {
	//TODO: see how accurate this is
	controls := make([]Vector2[types.Twip], 0, len(r.Control))
	for _, c := range r.Control {
		controls = append(controls, Vector2ToType[float64, types.Twip](transform.ApplyToVector(c.Float64().Divide(types.TwipFactor), applyTranslation).Multiply(types.TwipFactor)))
	}
	return &CubicSplineCurveRecord{
		Control: controls,
		Anchor:  Vector2ToType[float64, types.Twip](transform.ApplyToVector(r.Anchor.Float64().Divide(types.TwipFactor), applyTranslation).Multiply(types.TwipFactor)),
		Start:   Vector2ToType[float64, types.Twip](transform.ApplyToVector(r.Start.Float64().Divide(types.TwipFactor), applyTranslation).Multiply(types.TwipFactor)),
	}
}

func (r *CubicSplineCurveRecord) Equals(other Record) bool {
	if o, ok := other.(*CubicSplineCurveRecord); ok {
		return reflect.DeepEqual(r.Control, o.Control) && r.Start == o.Start && r.Anchor == o.Anchor
	}
	return false
}

func (r *CubicSplineCurveRecord) IsFlat() bool {
	return false
}
