package types

import (
	"golang.org/x/exp/maps"
	"slices"
)

type SpriteDefinition struct {
	ObjectId     uint16
	Processor    *SWFTreeProcessor
	CurrentFrame *ViewFrame
}

func (d *SpriteDefinition) GetObjectId() uint16 {
	return d.ObjectId
}

func (d *SpriteDefinition) GetShapeList(ratio float64) (list DrawPathList) {
	if d.CurrentFrame != nil {
		for _, object := range d.CurrentFrame.Render(0, nil, nil, nil) {
			list = append(list, object.DrawPathList...)
		}
	}
	return list
}

func (d *SpriteDefinition) NextFrame() *ViewFrame {
	//TODO: figure out why this can return null. missing shapes?
	d.CurrentFrame = d.NextFrame()
	if d.CurrentFrame == nil {
		return NewViewFrame(d.GetObjectId(), &DrawPathList{})
	}
	return d.CurrentFrame
}

func (d *SpriteDefinition) GetSafeObject() ObjectDefinition {
	return &SpriteDefinition{
		ObjectId:  d.ObjectId,
		Processor: NewSWFTreeProcessor(d.ObjectId, slices.Clone(d.Processor.Tags), maps.Clone(d.Processor.Objects)),
	}
}
