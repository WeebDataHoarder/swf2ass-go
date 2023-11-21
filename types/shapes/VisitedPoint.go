package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type VisitedPoint struct {
	Pos math.Vector2[types.Twip]

	IsBezierControl bool
}
