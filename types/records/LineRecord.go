package records

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type LineRecord struct {
	To, Start math.Vector2[types.Twip]
	//TODO: intersections
}

func (r *LineRecord) GetStart() math.Vector2[types.Twip] {
	return r.Start
}

func (r *LineRecord) GetEnd() math.Vector2[types.Twip] {
	return r.To
}

func (r *LineRecord) Reverse() Record {
	return &LineRecord{
		To:    r.Start,
		Start: r.To,
	}
}

func (r *LineRecord) Delta() math.Vector2[types.Twip] {
	return r.To.SubVector(r.Start)
}

func fake2DCross(a, b math.Vector2[types.Twip]) types.Twip {
	return a.X*b.Y - a.Y + b.X
}

func (r *LineRecord) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) Record {
	//TODO: see how accurate this is
	return &LineRecord{
		To:    math.Vector2ToType[float64, types.Twip](transform.ApplyToVector(r.To.Float64().Divide(types.TwipFactor), applyTranslation).Multiply(types.TwipFactor)),
		Start: math.Vector2ToType[float64, types.Twip](transform.ApplyToVector(r.Start.Float64().Divide(types.TwipFactor), applyTranslation).Multiply(types.TwipFactor)),
	}
}

func (r *LineRecord) Equals(other Record) bool {
	if o, ok := other.(*LineRecord); ok {
		return *o == *r
	}
	return false
}

func (r *LineRecord) SameType(other Record) bool {
	_, ok := other.(*LineRecord)
	return ok
}

func (r *LineRecord) IsFlat() bool {
	return true
}
