package types

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type SpriteDefinition struct {
	ObjectId     uint16
	Processor    *SWFTreeProcessor
	CurrentFrame *ViewFrame
}

func (d *SpriteDefinition) GetObjectId() uint16 {
	return d.ObjectId
}

func (d *SpriteDefinition) GetShapeList(p shapes.ObjectProperties) (list shapes.DrawPathList) {
	if d.CurrentFrame != nil {
		for _, object := range d.CurrentFrame.Render(0, nil, nil, nil) {
			list = append(list, object.DrawPathList...)
		}
	}
	panic("should not be called")
	return list
}

func (d *SpriteDefinition) NextFrame() *ViewFrame {
	//TODO: figure out why this can return null. missing shapes?
	d.CurrentFrame = d.Processor.NextFrame()
	if d.CurrentFrame == nil {
		return NewViewFrame(d.GetObjectId(), &shapes.DrawPathList{})
	}
	return d.CurrentFrame
}

func (d *SpriteDefinition) GetSafeObject() shapes.ObjectDefinition {
	return &SpriteDefinition{
		ObjectId:  d.ObjectId,
		Processor: NewSWFTreeProcessor(d.ObjectId, d.Processor.Tags, d.Processor.Objects),
	}
}
