package records

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type Record interface {
	GetStart() math.Vector2[float64]
	GetEnd() math.Vector2[float64]

	Reverse() Record

	Equals(other Record) bool

	SameType(other Record) bool

	ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) Record

	IsFlat() bool
}

type CurvedRecord interface {
	Record
	ToLineRecords(scale int64) []Record
}

type RecordPair [2]Record
