package records

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"reflect"
	"slices"
)

type CubicSplineCurveRecord struct {
	Control []math.Vector2[types.Twip]
	Anchor  math.Vector2[types.Twip]
	Start   math.Vector2[types.Twip]
}

func (r *CubicSplineCurveRecord) GetStart() math.Vector2[types.Twip] {
	return r.Start
}

func (r *CubicSplineCurveRecord) GetEnd() math.Vector2[types.Twip] {
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

func (r *CubicSplineCurveRecord) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) Record {
	//TODO: see how accurate this is
	controls := make([]math.Vector2[types.Twip], 0, len(r.Control))
	for _, c := range r.Control {
		controls = append(controls, math.Vector2ToType[float64, types.Twip](transform.ApplyToVector(c.Float64().Divide(types.TwipFactor), applyTranslation).Multiply(types.TwipFactor)))
	}
	return &CubicSplineCurveRecord{
		Control: controls,
		Anchor:  math.Vector2ToType[float64, types.Twip](transform.ApplyToVector(r.Anchor.Float64().Divide(types.TwipFactor), applyTranslation).Multiply(types.TwipFactor)),
		Start:   math.Vector2ToType[float64, types.Twip](transform.ApplyToVector(r.Start.Float64().Divide(types.TwipFactor), applyTranslation).Multiply(types.TwipFactor)),
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
