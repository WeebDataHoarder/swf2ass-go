package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/line"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type Tag interface {
	Equals(tag Tag) bool
	Encode(event time.EventTime) string
}

type StyleTag interface {
	Tag
	TransitionStyleRecord(line *line.Line, record shapes.StyleRecord) StyleTag
	FromStyleRecord(record shapes.StyleRecord) StyleTag
}

type PositioningTag interface {
	Tag
	TransitionMatrixTransform(line *line.Line, transform math.MatrixTransform) PositioningTag
	FromMatrixTransform(transform math.MatrixTransform) PositioningTag
}

type PathTag interface {
	Tag
	TransitionShape(line *line.Line, shape *shapes.Shape) PathTag
}

type ClipPathTag interface {
	Tag
	TransitionClipPath(line *line.Line, clip *types.ClipPath) ClipPathTag
}

type ColorTag interface {
	Tag
	ApplyColorTransform(transform math.ColorTransform) ColorTag
	TransitionColor(line *line.Line, transform math.ColorTransform) ColorTag
}
