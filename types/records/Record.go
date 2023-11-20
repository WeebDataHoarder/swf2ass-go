package records

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type Record interface {
	GetStart() math.Vector2[types.Twip]
	GetEnd() math.Vector2[types.Twip]

	Reverse() Record

	Equals(other Record) bool

	ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) Record

	IsFlat() bool
}

type RecordPair [2]Record
