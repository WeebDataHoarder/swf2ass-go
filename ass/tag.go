package ass

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
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
	TransitionMatrixTransform(line *Line, transform types.MatrixTransform) PositioningTag
	FromMatrixTransform(transform types.MatrixTransform) PositioningTag
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
	ApplyColorTransform(transform types.ColorTransform) ColorTag
	TransitionColor(line *Line, transform types.ColorTransform) ColorTag
}
