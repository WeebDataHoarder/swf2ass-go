package ass

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type Tag interface {
	Equals(tag Tag) bool
	Encode(event EventTime) string
}

type StyleTag interface {
	Tag
	TransitionStyleRecord(line *Line, record types.StyleRecord) StyleTag
	FromStyleRecord(record types.StyleRecord) StyleTag
}

type PositioningTag interface {
	Tag
	TransitionMatrixTransform(line *Line, transform math.MatrixTransform) PositioningTag
	FromMatrixTransform(transform math.MatrixTransform) PositioningTag
}

type PathTag interface {
	Tag
	TransitionShape(line *Line, shape *types.Shape) PathTag
}

type ClipPathTag interface {
	Tag
	TransitionClipPath(line *Line, clip *types.ClipPath) ClipPathTag
}

type ColorTag interface {
	Tag
	ApplyColorTransform(transform math.ColorTransform) ColorTag
	TransitionColor(line *Line, transform math.ColorTransform) ColorTag
}
