package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type Tag interface {
	Equals(tag Tag) bool
	Encode(event time.EventTime) string
}

type Event interface {
	GetStart() int64
	GetEnd() int64
}

type StyleTag interface {
	Tag
	TransitionStyleRecord(event Event, record shapes.StyleRecord) StyleTag
	FromStyleRecord(record shapes.StyleRecord) StyleTag
}

type PositioningTag interface {
	Tag
	TransitionMatrixTransform(event Event, transform math.MatrixTransform) PositioningTag
	FromMatrixTransform(transform math.MatrixTransform) PositioningTag
}

type PathTag interface {
	Tag
	TransitionShape(event Event, shape *shapes.Shape) PathTag
}

type ClipPathTag interface {
	Tag
	TransitionClipPath(event Event, clip *types.ClipPath) ClipPathTag
}

type ColorTag interface {
	Tag
	ApplyColorTransform(transform math.ColorTransform) ColorTag
	TransitionColor(event Event, transform math.ColorTransform) ColorTag
}
