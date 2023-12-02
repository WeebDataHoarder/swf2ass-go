package types

import (
	"fmt"
	math2 "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	"golang.org/x/exp/maps"
	"math"
	"slices"
)

type ViewLayout struct {
	Parent *ViewLayout

	DepthMap map[uint16]*ViewLayout

	Object shapes.ObjectDefinition

	ColorTransform  Option[math2.ColorTransform]
	MatrixTransform Option[math2.MatrixTransform]

	Properties shapes.ObjectProperties

	ClipDepth Option[uint16]
}

func NewClippingViewLayout(objectId, clipDepth uint16, object shapes.ObjectDefinition, parent *ViewLayout) *ViewLayout {
	l := NewViewLayout(objectId, object, parent)
	l.ClipDepth = Some(clipDepth)
	return l
}

func NewViewLayout(objectId uint16, object shapes.ObjectDefinition, parent *ViewLayout) *ViewLayout {
	if object != nil && object.GetObjectId() != objectId {
		panic("logic error")
	}
	return &ViewLayout{
		Parent:   parent,
		Object:   object,
		DepthMap: make(map[uint16]*ViewLayout),
		Properties: shapes.ObjectProperties{
			Visible: true,
		},
	}
}

func (v *ViewLayout) GetObjectId() uint16 {
	if v.Object != nil {
		return v.Object.GetObjectId()
	}
	return math.MaxUint16
}

func (v *ViewLayout) Get(depth uint16) *ViewLayout {
	return v.DepthMap[depth]
}

func (v *ViewLayout) Replace(depth uint16, ob *ViewLayout) {
	if v.Object != nil {
		panic("Cannot have ObjectDefinition and children at the same time")
	} else if oldObject, ok := v.DepthMap[depth]; ok && oldObject != nil {
		if _, ok := ob.MatrixTransform.Some(); !ok {
			ob.MatrixTransform = oldObject.MatrixTransform
		}
		if _, ok := ob.ColorTransform.Some(); !ok {
			ob.ColorTransform = oldObject.ColorTransform
		}
	}
	v.DepthMap[depth] = ob
}

func (v *ViewLayout) Place(depth uint16, ob *ViewLayout) {
	if v.Object != nil {
		panic("Cannot have ObjectDefinition and children at the same time")
	} else if ow, ok := v.DepthMap[depth]; ok && ow != ob {
		panic(fmt.Sprintf("Depth %d already exists: tried replacing object %d with %d", depth, ow.GetObjectId(), ob.GetObjectId()))
	}
	v.DepthMap[depth] = ob
}

func (v *ViewLayout) Remove(depth uint16) {
	delete(v.DepthMap, depth)
}

func (v *ViewLayout) NextFrame(frameNumber int64, actions ActionList) (frame *ViewFrame) {
	frame = v.nextFrame(frameNumber, actions)

	if clipDepth, isClipping := v.ClipDepth.Some(); isClipping {
		clip := NewClippingFrame(frame.ObjectId, clipDepth, frame.DrawPathList)
		for depth, f := range frame.DepthMap {
			clip.AddChild(depth, f)
		}
		clip.ColorTransform = frame.ColorTransform
		clip.MatrixTransform = frame.MatrixTransform
		return clip
	}
	return frame
}

func (v *ViewLayout) nextFrame(frameNumber int64, actions ActionList) (frame *ViewFrame) {
	if v.Object != nil {
		if mfod, ok := v.Object.(MultiFrameObjectDefinition); ok {
			frame = mfod.NextFrame(frameNumber, v.Properties)
		} else {
			list := v.Object.GetShapeList(v.Properties)
			frame = NewViewFrame(v.GetObjectId(), Some(list))
		}
	} else {
		frame = NewViewFrame(v.GetObjectId(), None[shapes.DrawPathList]())

		keys := maps.Keys(v.DepthMap)
		slices.Sort(keys)

		for _, depth := range keys {
			child := v.DepthMap[depth]
			f := child.NextFrame(frameNumber, actions)
			frame.AddChild(depth, f)
		}
	}

	frame.ColorTransform = v.ColorTransform
	frame.MatrixTransform = v.MatrixTransform

	return frame
}
