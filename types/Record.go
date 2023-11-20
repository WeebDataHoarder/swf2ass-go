package types

import "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"

type Record interface {
	GetStart() Vector2[types.Twip]
	GetEnd() Vector2[types.Twip]

	Reverse() Record

	Equals(other Record) bool

	ApplyMatrixTransform(transform MatrixTransform, applyTranslation bool) Record

	IsFlat() bool
}
