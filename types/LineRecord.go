package types

import "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"

type LineRecord struct {
	To, Start Vector2[types.Twip]
	//TODO: intersections
}

func (r *LineRecord) GetStart() Vector2[types.Twip] {
	return r.Start
}

func (r *LineRecord) GetEnd() Vector2[types.Twip] {
	return r.To
}

func (r *LineRecord) Reverse() Record {
	return &LineRecord{
		To:    r.Start,
		Start: r.To,
	}
}

func (r *LineRecord) ApplyMatrixTransform(transform MatrixTransform, applyTranslation bool) Record {
	//TODO: see how accurate this is
	return &LineRecord{
		To:    Vector2ToType[float64, types.Twip](transform.ApplyToVector(r.To.Float64().Divide(types.TwipFactor), applyTranslation).Multiply(types.TwipFactor)),
		Start: Vector2ToType[float64, types.Twip](transform.ApplyToVector(r.Start.Float64().Divide(types.TwipFactor), applyTranslation).Multiply(types.TwipFactor)),
	}
}

func (r *LineRecord) Equals(other Record) bool {
	if o, ok := other.(*LineRecord); ok {
		return *o == *r
	}
	return false
}

func (r *LineRecord) IsFlat() bool {
	return true
}
