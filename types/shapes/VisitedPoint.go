package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type VisitedPoint[T ~float64 | ~int64] struct {
	Pos math.Vector2[T]

	IsBezierControl bool
}
