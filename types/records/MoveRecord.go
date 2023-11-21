package records

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type MoveRecord struct {
	To, Start math.Vector2[types.Twip]
}

func (r *MoveRecord) GetStart() math.Vector2[types.Twip] {
	return r.Start
}

func (r *MoveRecord) GetEnd() math.Vector2[types.Twip] {
	return r.To
}

func (r *MoveRecord) Reverse() Record {
	return &MoveRecord{
		To:    r.Start,
		Start: r.To,
	}
}

func (r *MoveRecord) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) Record {
	//TODO: see how accurate this is
	return &MoveRecord{
		To:    math.Vector2ToType[float64, types.Twip](transform.ApplyToVector(r.To.Float64().Divide(types.TwipFactor), applyTranslation).Multiply(types.TwipFactor)),
		Start: math.Vector2ToType[float64, types.Twip](transform.ApplyToVector(r.Start.Float64().Divide(types.TwipFactor), applyTranslation).Multiply(types.TwipFactor)),
	}
}

func (r *MoveRecord) Equals(other Record) bool {
	if o, ok := other.(*MoveRecord); ok {
		return *o == *r
	}
	return false
}

func (r *MoveRecord) IsFlat() bool {
	return true
}