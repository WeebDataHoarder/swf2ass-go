package types

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	math2 "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
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

func (d *MorphShapeDefinition) GetShapeList(ratio float64) (list shapes.DrawPathList) {
	//TODO: cache shapes by ratio
	//TODO: refactor this to use color transforms (and if able) matrix transforms

	if math.Abs(ratio) < math.SmallestNonzeroFloat64 {
		return d.StartShapeList
	}
	if math.Abs(ratio-1.0) < math.SmallestNonzeroFloat64 {
		return d.EndShapeList
	}

	for i, c1 := range d.StartShapeList {
		c2 := d.EndShapeList[i]

		var shape shapes.Shape

		for _, recordPair := range shapes.IterateMorphShape(c1.Commands, c2.Commands) {
			startEdge := recordPair[0]
			endEdge := recordPair[1]

			//No need to convert types!
			aLineRecord, aIsLineRecord := startEdge.(*records.LineRecord)
			aMoveRecord, aIsMoveRecord := startEdge.(*records.MoveRecord)
			aQuadraticCurveRecord, aIsQuadraticCurveRecord := startEdge.(*records.QuadraticCurveRecord)
			bLineRecord, bIsLineRecord := endEdge.(*records.LineRecord)
			bMoveRecord, bIsMoveRecord := endEdge.(*records.MoveRecord)
			bQuadraticCurveRecord, bIsQuadraticCurveRecord := endEdge.(*records.QuadraticCurveRecord)

			if aIsLineRecord && bIsLineRecord {
				shape.AddRecord(&records.LineRecord{
					To:    math2.LerpVector2(aLineRecord.To, bLineRecord.To, ratio),
					Start: math2.LerpVector2(aLineRecord.Start, bLineRecord.Start, ratio),
				})
			} else if aIsQuadraticCurveRecord && bIsQuadraticCurveRecord {
				shape.AddRecord(&records.QuadraticCurveRecord{
					Control: math2.LerpVector2(aQuadraticCurveRecord.Control, bQuadraticCurveRecord.Control, ratio),
					Anchor:  math2.LerpVector2(aQuadraticCurveRecord.Anchor, bQuadraticCurveRecord.Anchor, ratio),
					Start:   math2.LerpVector2(aQuadraticCurveRecord.Start, bQuadraticCurveRecord.Start, ratio),
				})
			} else if aIsMoveRecord && bIsMoveRecord {
				shape.AddRecord(&records.MoveRecord{
					To:    math2.LerpVector2(aMoveRecord.To, bMoveRecord.To, ratio),
					Start: math2.LerpVector2(aMoveRecord.Start, bMoveRecord.Start, ratio),
				})
			} else {
				panic("unsupported")
			}
		}

		//TODO: morph styles properly
		c1FillStyle, c1IsFillStyle := c1.Style.(*shapes.FillStyleRecord)
		c1LineStyle, c1IsLineStyle := c1.Style.(*shapes.LineStyleRecord)
		c2FillStyle, c2IsFillStyle := c2.Style.(*shapes.FillStyleRecord)
		c2LineStyle, c2IsLineStyle := c2.Style.(*shapes.LineStyleRecord)

		if c1IsFillStyle && c2IsFillStyle {
			if c1Color, ok := c1FillStyle.Fill.(math2.Color); ok {
				list = append(list, shapes.DrawPathFill(&shapes.FillStyleRecord{
					Fill:   math2.LerpColor(c1Color, c2FillStyle.Fill.(math2.Color), ratio),
					Border: c1FillStyle.Border,
				}, &shape))
			} else if c1Gradient, ok := c1FillStyle.Fill.(shapes.Gradient); ok {
				//TODO: proper gradients
				list = append(list, shapes.DrawPathFill(&shapes.FillStyleRecord{
					Fill:   math2.LerpColor(c1Gradient.GetItems()[0].Color, c2FillStyle.Fill.(shapes.Gradient).GetItems()[0].Color, ratio),
					Border: c1FillStyle.Border,
				}, &shape))
			} else {
				panic("unsupported")
			}
		} else if c1IsLineStyle && c2IsLineStyle {
			list = append(list, shapes.DrawPathStroke(&shapes.LineStyleRecord{
				Width: math2.Lerp(c1LineStyle.Width, c2LineStyle.Width, ratio),
				Color: math2.LerpColor(c1LineStyle.Color, c2LineStyle.Color, ratio),
			}, &shape))
		} else {
			panic("unsupported")
		}
	}

	return list
}

func (d *MorphShapeDefinition) GetSafeObject() ObjectDefinition {
	return d
}

func MorphShapeDefinitionFromSWF(shapeId uint16, startBounds, endBounds shapes.Rectangle[float64], startRecords, endRecords subtypes.SHAPERECORDS, fillStyles subtypes.MORPHFILLSTYLEARRAY, lineStyles subtypes.MORPHLINESTYLEARRAY) *MorphShapeDefinition {
	startStyles, endStyles := shapes.StyleListFromSWFMorphItems(fillStyles, lineStyles)

	start := shapes.DrawPathListFromSWFMorph(startRecords, endRecords, startStyles, false)
	//TODO: morph styles properly
	_ = endStyles
	end := shapes.DrawPathListFromSWFMorph(startRecords, endRecords, startStyles, true)

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
