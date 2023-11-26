package types

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/records"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	"math"
)

type MorphShapeDefinition struct {
	ObjectId                     uint16
	StartBounds, EndBounds       shapes.Rectangle[float64]
	StartShapeList, EndShapeList shapes.DrawPathList
}

func (d *MorphShapeDefinition) GetObjectId() uint16 {
	return d.ObjectId
}

func (d *MorphShapeDefinition) GetShapeList(p shapes.ObjectProperties) (list shapes.DrawPathList) {
	//TODO: cache shapes by ratio
	//TODO: refactor this to use color transforms (and if able) matrix transforms

	if math.Abs(p.Ratio) < math.SmallestNonzeroFloat64 {
		return d.StartShapeList
	}
	if math.Abs(p.Ratio-1.0) < math.SmallestNonzeroFloat64 {
		return d.EndShapeList
	}

	for i, c1 := range d.StartShapeList {
		c2 := d.EndShapeList[i]

		var shape shapes.Shape

		for _, recordPair := range shapes.IterateMorphShape(c1.Shape, c2.Shape) {
			shape = append(shape, records.LerpRecord(recordPair[0], recordPair[1], p.Ratio))
		}

		//TODO: morph styles properly
		c1FillStyle, c1IsFillStyle := c1.Style.(*shapes.FillStyleRecord)
		c1LineStyle, c1IsLineStyle := c1.Style.(*shapes.LineStyleRecord)
		c2FillStyle, c2IsFillStyle := c2.Style.(*shapes.FillStyleRecord)
		c2LineStyle, c2IsLineStyle := c2.Style.(*shapes.LineStyleRecord)

		if c1IsFillStyle && c2IsFillStyle {
			list = append(list, shapes.DrawPathFill(shapes.LerpFillStyle(c1FillStyle, c2FillStyle, p.Ratio), shape))
		} else if c1IsLineStyle && c2IsLineStyle {
			list = append(list, shapes.DrawPathStroke(shapes.LerpLineStyle(c1LineStyle, c2LineStyle, p.Ratio), shape))
		} else {
			panic("unsupported")
		}
	}

	return list
}

func (d *MorphShapeDefinition) GetSafeObject() shapes.ObjectDefinition {
	return d
}

func MorphShapeDefinitionFromSWF(collection shapes.ObjectCollection, shapeId uint16, startBounds, endBounds shapes.Rectangle[float64], startRecords, endRecords subtypes.SHAPERECORDS, fillStyles subtypes.MORPHFILLSTYLEARRAY, lineStyles subtypes.MORPHLINESTYLEARRAY) *MorphShapeDefinition {
	startStyles, endStyles := shapes.StyleListFromSWFMorphItems(collection, fillStyles, lineStyles)

	start, end := shapes.DrawPathListFromSWFMorph(collection, startRecords, endRecords, startStyles, endStyles)

	if len(start) != len(end) {
		panic("length does not match")
	}

	return &MorphShapeDefinition{
		ObjectId:       shapeId,
		StartBounds:    startBounds,
		EndBounds:      endBounds,
		StartShapeList: start,
		EndShapeList:   end,
	}
}
