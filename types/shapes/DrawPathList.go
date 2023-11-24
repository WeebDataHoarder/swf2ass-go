package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type DrawPathList []DrawPath

func (l DrawPathList) Merge(b DrawPathList) DrawPathList {
	newList := make(DrawPathList, 0, len(l)+len(b))
	newList = append(newList, l...)
	newList = append(newList, b...)
	return newList
}

func (l DrawPathList) ApplyFunction(f func(p DrawPath) DrawPath) (r DrawPathList) {
	r = make(DrawPathList, 0, len(l))
	for _, e := range l {
		r = append(r, f(e))
	}
	return r
}

func (l DrawPathList) Fill(shape *Shape) (r DrawPathList) {
	clipShape := NewClipPath(shape)
	//Convert paths to many tags using intersections
	for _, innerPath := range l {
		newPath := DrawPath{
			Style:    innerPath.Style,
			Commands: innerPath.Commands,
			Clip:     clipShape,
		}
		if len(newPath.Commands.Edges) == 0 {
			continue
		}

		r = append(r, newPath)
	}
	return r
}

func (l DrawPathList) ApplyColorTransform(transform math.ColorTransform) Fillable {
	r := make(DrawPathList, 0, len(l))
	for _, e := range l {
		r = append(r, DrawPath{
			Style:    e.Style.ApplyColorTransform(transform),
			Commands: e.Commands,
		})
	}
	return r
}

func (l DrawPathList) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) (r DrawPathList) {
	r = make(DrawPathList, 0, len(l))
	for i := range l {
		r = append(r, l[i].ApplyMatrixTransform(transform, applyTranslation))
	}
	return r
}

func DrawPathListFromSWF(collection ObjectCollection, records subtypes.SHAPERECORDS, styles StyleList) DrawPathList {
	converter := NewShapeConverter(collection, records, styles)
	converter.Convert(false)

	return converter.Commands
}

func DrawPathListFromSWFMorph(collection ObjectCollection, startRecords, endRecords subtypes.SHAPERECORDS, styles StyleList, flip bool) DrawPathList {
	converter := NewMorphShapeConverter(collection, startRecords, endRecords, styles)
	converter.Convert(flip)

	return converter.Commands
}
