package types

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	"image"
)

type BitmapDefinition struct {
	ObjectId  uint16
	ShapeList shapes.DrawPathList
}

func (d *BitmapDefinition) GetObjectId() uint16 {
	return d.ObjectId
}

func (d *BitmapDefinition) GetShapeList(ratio float64) (list shapes.DrawPathList) {
	return d.ShapeList
}

func (d *BitmapDefinition) GetSafeObject() shapes.ObjectDefinition {
	return d
}

func BitmapDefinitionFromSWF(bitmapId uint16, imageData []byte, alphaData []byte) (*BitmapDefinition, error) {
	l, err := shapes.ConvertBitmapBytesToDrawPathList(imageData, alphaData)
	if err != nil {
		return nil, err
	}

	return &BitmapDefinition{
		ObjectId:  bitmapId,
		ShapeList: l,
	}, nil
}

func BitmapDefinitionFromSWFLossless(bitmapId uint16, im image.Image) *BitmapDefinition {
	if im == nil {
		return nil
	}
	return &BitmapDefinition{
		ObjectId:  bitmapId,
		ShapeList: shapes.ConvertBitmapToDrawPathList(im),
	}
}
