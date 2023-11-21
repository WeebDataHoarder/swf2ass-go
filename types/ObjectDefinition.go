package types

import "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"

type ObjectDefinition interface {
	GetObjectId() uint16
	GetSafeObject() ObjectDefinition
	GetShapeList(ratio float64) shapes.DrawPathList
}

type MultiFrameObjectDefinition interface {
	ObjectDefinition
	NextFrame() *ViewFrame
}
