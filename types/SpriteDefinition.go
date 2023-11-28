package types

import (
	math2 "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	"slices"
)

type SpriteDefinition struct {
	ObjectId uint16
	Frames   [][]SpriteFrameEntry
}

func (d *SpriteDefinition) GetObjectId() uint16 {
	return d.ObjectId
}

func (d *SpriteDefinition) GetShapeList(p shapes.ObjectProperties) (list shapes.DrawPathList) {
	panic("should not be called")
	return list
}

func (d *SpriteDefinition) NextFrame(frameNumber int64, p shapes.ObjectProperties) *ViewFrame {
	frameN := frameNumber - p.PlaceFrame

	n := frameN % int64(len(d.Frames))

	spriteFrame := NewViewFrame(d.ObjectId, nil)

	for _, e := range d.Frames[n] {
		var frame *ViewFrame
		if mfod, ok := e.Object.(MultiFrameObjectDefinition); ok {
			frame = mfod.NextFrame(frameN, e.Properties)
		} else {
			list := e.Object.GetShapeList(e.Properties)
			frame = NewViewFrame(e.Object.GetObjectId(), &list)
		}

		frame.ColorTransform = e.ColorTransform
		frame.MatrixTransform = e.MatrixTransform

		frame.ClipDepth = e.ClipDepth

		spriteFrame.AddChild(e.Depth, frame)
	}

	return spriteFrame
}

func (d *SpriteDefinition) GetSafeObject() shapes.ObjectDefinition {
	return d
}

type SpriteFrameEntry struct {
	Depth           uint16
	Object          shapes.ObjectDefinition
	ColorTransform  Option[math2.ColorTransform]
	MatrixTransform Option[math2.MatrixTransform]
	ClipDepth       Option[uint16]
	Properties      shapes.ObjectProperties
}

func SpriteDefinitionFromSWF(spriteId uint16, frameCount int, p *SWFTreeProcessor) *SpriteDefinition {
	var frames [][]SpriteFrameEntry

	var lastFrame *ViewFrame
	for p.Loops == 0 && (len(frames) < frameCount || (frameCount == 0 && len(frames) == 0)) {
		f := p.NextFrame()
		if f == nil {
			break
		}
		if lastFrame == f {
			break
		}

		lastFrame = f

		var entries []SpriteFrameEntry

		for depth, layout := range p.Layout.DepthMap {
			if layout.Object == nil {
				panic("not supported")
			}

			entries = append(entries, SpriteFrameEntry{
				Depth:           depth,
				Object:          layout.Object,
				ColorTransform:  layout.ColorTransform,
				MatrixTransform: layout.MatrixTransform,
				ClipDepth:       layout.ClipDepth,
				Properties:      layout.Properties,
			})
		}
		slices.SortFunc(entries, func(a, b SpriteFrameEntry) int {
			return int(a.Depth) - int(b.Depth)
		})
		frames = append(frames, entries)
	}

	if len(frames) == 0 {
		panic("unsupported")
	}

	return &SpriteDefinition{
		ObjectId: spriteId,
		Frames:   frames,
	}
}
