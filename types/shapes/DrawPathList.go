package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
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
	if false { //TODO

		clipShape := NewClipPath(shape)
		//Convert paths to many tags using intersections
		for _, innerPath := range l {
			newPath := DrawPath{
				Style:    innerPath.Style,
				Commands: clipShape.ClipShape(innerPath.Commands),
			}
			if len(newPath.Commands.Edges) == 0 {
				continue
			}

			r = append(r, newPath)
		}
		return r
	}

	//TODO: fix this below
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
	for i := range l {
		r = append(r, l[i].ApplyColorTransform(transform))
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

func DrawPathListFillFromSWF(l DrawPathList, transform types.MATRIX) DrawPathList {
	// shape is already in pixel world, but matrix comes as twip
	baseScale := math.ScaleTransform(math.NewVector2[float64](1./types.TwipFactor, 1./types.TwipFactor))
	t := math.MatrixTransformFromSWF(transform).Multiply(baseScale)
	return l.ApplyMatrixTransform(t, true)
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
