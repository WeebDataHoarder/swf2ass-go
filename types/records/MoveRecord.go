package records

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type MoveRecord struct {
	To, Start math.Vector2[float64]
}

func (r MoveRecord) GetStart() math.Vector2[float64] {
	return r.Start
}

func (r MoveRecord) GetEnd() math.Vector2[float64] {
	return r.To
}

func (r MoveRecord) Reverse() Record {
	return MoveRecord{
		To:    r.Start,
		Start: r.To,
	}
}

func (r MoveRecord) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) Record {
	//TODO: see how accurate this is
	return MoveRecord{
		To:    math.MatrixTransformApplyToVector(transform, r.To, applyTranslation),
		Start: math.MatrixTransformApplyToVector(transform, r.Start, applyTranslation),
	}
}

func (r MoveRecord) Equals(other Record) bool {
	if o, ok := other.(MoveRecord); ok {
		return o == r
	}
	return false
}

func (r MoveRecord) SameType(other Record) bool {
	_, ok := other.(MoveRecord)
	return ok
}

func (r MoveRecord) IsFlat() bool {
	return true
}
