package records

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type LineRecord struct {
	To, Start math.Vector2[float64]
	//TODO: intersections
}

func (r LineRecord) GetStart() math.Vector2[float64] {
	return r.Start
}

func (r LineRecord) GetEnd() math.Vector2[float64] {
	return r.To
}

func (r LineRecord) Reverse() Record {
	return LineRecord{
		To:    r.Start,
		Start: r.To,
	}
}

func (r LineRecord) Delta() math.Vector2[float64] {
	return r.To.SubVector(r.Start)
}

func (r LineRecord) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) Record {
	//TODO: see how accurate this is
	return LineRecord{
		To:    transform.ApplyToVector(r.To, applyTranslation),
		Start: transform.ApplyToVector(r.Start, applyTranslation),
	}
}

func (r LineRecord) Equals(other Record) bool {
	if o, ok := other.(LineRecord); ok {
		return o == r
	}
	return false
}

func (r LineRecord) SameType(other Record) bool {
	_, ok := other.(LineRecord)
	return ok
}

func (r LineRecord) BoundingBox() (topLeft, bottomRight math.Vector2[float64]) {
	return r.Start.Min(r.To), r.Start.Max(r.To)
}

func (r LineRecord) IsFlat() bool {
	return true
}

func (r LineRecord) String() string {
	return fmt.Sprintf("l %s", r.To)
}
