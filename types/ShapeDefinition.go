package types

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type ShapeDefinition struct {
	ObjectId  uint16
	Bounds    shapes.Rectangle[float64]
	ShapeList shapes.DrawPathList
}

func (d *ShapeDefinition) GetObjectId() uint16 {
	return d.ObjectId
}

func (d *ShapeDefinition) GetShapeList(ratio float64) (list shapes.DrawPathList) {
	return d.ShapeList
}

func (d *ShapeDefinition) GetSafeObject() shapes.ObjectDefinition {
	return d
}

func ShapeDefinitionFromSWF(collection shapes.ObjectCollection, shapeId uint16, bounds shapes.Rectangle[float64], records subtypes.SHAPERECORDS, fillStyles subtypes.FILLSTYLEARRAY, lineStyles subtypes.LINESTYLEARRAY) *ShapeDefinition {
	styles := shapes.StyleListFromSWFItems(collection, fillStyles, lineStyles)

	drawPathList := shapes.DrawPathListFromSWF(collection, records, styles)

	return &ShapeDefinition{
		ObjectId:  shapeId,
		Bounds:    bounds,
		ShapeList: drawPathList,
	}
}
