package types

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type ShapeDefinition struct {
	ObjectId  uint16
	Bounds    shapes.Rectangle[types.Twip]
	ShapeList DrawPathList
}

func (d *ShapeDefinition) GetObjectId() uint16 {
	return d.ObjectId
}

func (d *ShapeDefinition) GetShapeList(ratio float64) (list DrawPathList) {
	return d.ShapeList
}

func (d *ShapeDefinition) GetSafeObject() ObjectDefinition {
	return d
}
