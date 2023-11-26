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

	ColorTransform  *math2.ColorTransform
	MatrixTransform *math2.MatrixTransform

	Properties shapes.ObjectProperties

	IsClipping bool
	ClipDepth  uint16
}

func NewClippingViewLayout(objectId, clipDepth uint16, object shapes.ObjectDefinition, parent *ViewLayout) *ViewLayout {
	l := NewViewLayout(objectId, object, parent)
	l.IsClipping = true
	l.ClipDepth = clipDepth
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
		if ob.MatrixTransform == nil {
			ob.MatrixTransform = oldObject.MatrixTransform
		}
		if ob.ColorTransform == nil {
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

func (v *ViewLayout) NextFrame(actions ActionList) (frame *ViewFrame) {
	frame = v.nextFrame(actions)

	if v.IsClipping {
		clip := NewClippingFrame(frame.ObjectId, v.ClipDepth, frame.DrawPathList)
		for depth, f := range frame.DepthMap {
			clip.AddChild(depth, f)
		}
		clip.ColorTransform = frame.ColorTransform
		clip.MatrixTransform = frame.MatrixTransform
		return clip
	}
	return frame
}

func (v *ViewLayout) nextFrame(actions ActionList) (frame *ViewFrame) {
	if v.Object != nil {
		if mfod, ok := v.Object.(MultiFrameObjectDefinition); ok {
			frame = mfod.NextFrame()
		} else {
			list := v.Object.GetShapeList(v.Properties)
			frame = NewViewFrame(v.GetObjectId(), &list)
		}
	} else {
		frame = NewViewFrame(v.GetObjectId(), nil)

		keys := maps.Keys(v.DepthMap)
		slices.Sort(keys)

		for _, depth := range keys {
			child := v.DepthMap[depth]
			f := child.NextFrame(actions)
			frame.AddChild(depth, f)
		}
	}

	frame.ColorTransform = v.ColorTransform
	frame.MatrixTransform = v.MatrixTransform

	return frame
}
