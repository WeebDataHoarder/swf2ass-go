package types

import "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"

type MultiFrameObjectDefinition interface {
	shapes.ObjectDefinition
	NextFrame(frameNumber int64, p shapes.ObjectProperties) *ViewFrame
}
