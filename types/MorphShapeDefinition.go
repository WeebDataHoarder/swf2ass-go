package types

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type MorphShapeDefinition struct {
	ObjectId                     uint16
	StartBounds, EndBounds       shapes.Rectangle[types.Twip]
	StartShapeList, EndShapeList shapes.DrawPathList
}

func (d *MorphShapeDefinition) GetObjectId() uint16 {
	return d.ObjectId
}

func (d *MorphShapeDefinition) GetShapeList(ratio float64) (list shapes.DrawPathList) {
	//TODO: implement morphs
	return d.StartShapeList
}

func (d *MorphShapeDefinition) GetSafeObject() ObjectDefinition {
	return d
}

func MorphShapeDefinitionFromSWF(shapeId uint16, startBounds, endBounds shapes.Rectangle[types.Twip], startRecords, endRecords subtypes.SHAPERECORDS, fillStyles subtypes.MORPHFILLSTYLEARRAY, lineStyles subtypes.MORPHLINESTYLEARRAY) *MorphShapeDefinition {
	startStyles, endStyles := shapes.StyleListFromSWFMorphItems(fillStyles, lineStyles)

	start := shapes.DrawPathListFromSWFMorph(startRecords, endRecords, startStyles, false)
	//TODO: morph styles properly
	_ = endStyles
	end := shapes.DrawPathListFromSWFMorph(startRecords, endRecords, startStyles, true)

	return &MorphShapeDefinition{
		ObjectId:       shapeId,
		StartBounds:    startBounds,
		EndBounds:      endBounds,
		StartShapeList: start,
		EndShapeList:   end,
	}
}
