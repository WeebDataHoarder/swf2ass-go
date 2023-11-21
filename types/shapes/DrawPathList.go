package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
)

type DrawPathList []DrawPath

func (l DrawPathList) Merge(b DrawPathList) DrawPathList {
	newList := make(DrawPathList, 0, len(l)+len(b))
	newList = append(newList, l...)
	newList = append(newList, b...)
	return newList
}

func DrawPathListFromSWF(records subtypes.SHAPERECORDS, styles StyleList) DrawPathList {
	converter := NewShapeConverter(records, styles)
	converter.Convert(false)

	return converter.Commands
}

func DrawPathListFromSWFMorph(startRecords, endRecords subtypes.SHAPERECORDS, styles StyleList, flip bool) DrawPathList {
	converter := NewMorphShapeConverter(startRecords, endRecords, styles)
	converter.Convert(flip)

	return converter.Commands
}
