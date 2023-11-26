package types

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type FontDefinition struct {
	FontId uint16

	Name   []byte
	Italic bool
	Bold   bool
	Scale  float64

	Entries []FontEntry
}

func (d *FontDefinition) GetObjectId() uint16 {
	return d.FontId
}

func (d *FontDefinition) GetShapeList(p shapes.ObjectProperties) (list shapes.DrawPathList) {
	fmt.Printf("something is trying to place a Font as a character!!!!!!")
	return nil
}

func (d *FontDefinition) ComposeTextSWF(entries []subtypes.GLYPHENTRY, height, yOffset, xOffset float64, c math.Color) (list shapes.DrawPathList, newXOffset float64) {
	for _, g := range entries {
		e := d.Entries[g.Index]

		t := math.TranslateTransform(math.NewVector2(xOffset, yOffset)).Multiply(math.ScaleTransform(math.NewVector2(height/d.Scale, height/d.Scale)))
		for _, dp := range e.List.ApplyMatrixTransform(t, true) {
			if _, ok := dp.Style.(*shapes.FillStyleRecord); ok {
				list = append(list, shapes.DrawPathFill(&shapes.FillStyleRecord{
					Fill:   c,
					Border: nil,
					Blur:   0,
				}, dp.Commands, nil))
			} else {
				continue
				//skip these, TODO: why?
				list = append(list, shapes.DrawPathFill(&shapes.FillStyleRecord{
					Fill:   c,
					Border: nil,
					Blur:   0,
				}, dp.Commands, nil))
			}
		}
		xOffset += types.Twip(g.Advance).Float64()
	}
	return list, xOffset
}

func (d *FontDefinition) GetSafeObject() shapes.ObjectDefinition {
	return d
}

type FontEntry struct {
	Bounds shapes.Rectangle[float64]
	Code   uint32
	Offset uint32
	List   shapes.DrawPathList
}

type TextDefinition struct {
	ObjectId  uint16
	Bounds    shapes.Rectangle[float64]
	ShapeList shapes.DrawPathList
}

func (d *TextDefinition) GetObjectId() uint16 {
	return d.ObjectId
}

func (d *TextDefinition) GetShapeList(p shapes.ObjectProperties) (list shapes.DrawPathList) {
	return d.ShapeList
}

func (d *TextDefinition) GetSafeObject() shapes.ObjectDefinition {
	return d
}

func TextDefinitionFromSWF(collection shapes.ObjectCollection, characterId uint16, bounds types.RECT, textRecords subtypes.TEXTRECORDS, matrix types.MATRIX) *TextDefinition {
	var currentFont *FontDefinition
	var textHeight float64
	var xOffset, yOffset float64
	var r, l shapes.DrawPathList

	currentColor := math.Color{
		R:     0,
		G:     0,
		B:     0,
		Alpha: 255,
	}

	characterBounds := shapes.RectangleFromSWF(bounds)

	for _, g := range textRecords {
		if g.Flag.HasFont {
			ob := collection.Get(g.FontId)
			if ob == nil {
				//font not found
				return nil
			}
			if fd, ok := ob.(*FontDefinition); ok {
				currentFont = fd
			} else {
				//font not valid
				return nil
			}
			textHeight = types.Twip(g.TextHeight).Float64()
		}

		if g.Flag.HasColor {
			currentColor.R = g.Color.R()
			currentColor.G = g.Color.G()
			currentColor.B = g.Color.B()
			currentColor.Alpha = g.Color.A()
		}

		if g.Flag.HasXOffset {
			xOffset = types.Twip(g.XOffset).Float64()
		}

		if g.Flag.HasYOffset {
			yOffset = types.Twip(g.YOffset).Float64()
		}

		if currentFont == nil {
			// no font defined
			return nil
		}

		l, xOffset = currentFont.ComposeTextSWF(g.GlyphEntries, textHeight, yOffset, xOffset, currentColor)
		r = append(r, l...)
	}

	return &TextDefinition{
		ObjectId:  characterId,
		Bounds:    characterBounds,
		ShapeList: r.ApplyMatrixTransform(math.MatrixTransformFromSWF(matrix), true),
	}
}

func FontDefinitionFromSWF[T1 uint16 | uint32, T2 uint8 | uint16](fontId uint16, name []byte, hasLayout, hasCodeTable, italic, bold bool, boundsTable []types.RECT, shapeTable []subtypes.SHAPE, offsetTable []T1, codeTable []T2, scale float64) *FontDefinition {

	styleList := shapes.StyleList{
		//TODO: why is this needed????
		FillStyles: []*shapes.FillStyleRecord{{}},
		//TODO: why is this needed????
		LineStyles: []*shapes.LineStyleRecord{{}},
	}

	var entries []FontEntry
	for i, s := range shapeTable {
		var bounds = shapes.Rectangle[float64]{}
		if hasLayout {
			bounds = shapes.RectangleFromSWF(boundsTable[i])
		}
		var code uint32
		if hasCodeTable {
			code = uint32(codeTable[i])
		}
		drawPathList := shapes.DrawPathListFromSWF(nil, s.Records, styleList)
		entries = append(entries, FontEntry{
			Bounds: bounds,
			Offset: uint32(offsetTable[i]),
			Code:   code,
			List:   drawPathList,
		})
	}

	return &FontDefinition{
		FontId:  fontId,
		Italic:  italic,
		Scale:   scale / types.TwipFactor,
		Bold:    bold,
		Name:    name,
		Entries: entries,
	}
}
