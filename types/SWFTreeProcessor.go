package types

import (
	"fmt"
	swftag "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	math2 "math"
	"runtime"
)

type SWFTreeProcessor struct {
	Layout *ViewLayout

	Objects ObjectCollection

	Tags []swftag.Tag

	Index int

	Frame int64

	LastFrame *ViewFrame
	Playing   bool
	Loops     int

	processFunc func(actions ActionList) (tag swftag.Tag, newActions ActionList)
}

func NewSWFTreeProcessor(objectId uint16, tags []swftag.Tag, objects ObjectCollection) *SWFTreeProcessor {
	return &SWFTreeProcessor{
		Objects: objects,
		Frame:   0,
		Tags:    tags,
		Layout:  NewViewLayout(objectId, nil, nil),
		Playing: true,
	}
}

func (p *SWFTreeProcessor) Next() {
	p.Index++
}

func (p *SWFTreeProcessor) Current() swftag.Tag {
	if len(p.Tags) > p.Index {
		return p.Tags[p.Index]
	}
	return nil
}

func (p *SWFTreeProcessor) Process(actions ActionList) (tag swftag.Tag, newActions ActionList) {
	if p.processFunc != nil {
		return p.processFunc(actions)
	}
	return p.process(actions)
}

func (p *SWFTreeProcessor) placeObject(object ObjectDefinition, depth, clipDepth uint16, isMove, hasRatio, hasClipDepth bool, ratio float64, transform *math.MatrixTransform, colorTransform *math.ColorTransform) {
	if object == nil {
		//TODO: place bogus element
		fmt.Printf("Object at depth:%d not found\n", depth)
		p.Layout.Remove(depth)
		return
	}

	currentLayout := p.Layout.Get(depth)

	if isMove && currentLayout != nil && currentLayout.GetObjectId() == object.GetObjectId() {
		if transform != nil {
			currentLayout.MatrixTransform = transform
		}
		if colorTransform != nil {
			currentLayout.ColorTransform = colorTransform
		}
		if hasRatio {
			currentLayout.Ratio = ratio
		}
		return
	}

	var view *ViewLayout
	if hasClipDepth {
		view = NewClippingViewLayout(object.GetObjectId(), clipDepth, object.GetSafeObject(), p.Layout)
	} else {
		view = NewViewLayout(object.GetObjectId(), object.GetSafeObject(), p.Layout)
	}
	view.MatrixTransform = transform
	view.ColorTransform = colorTransform
	view.Ratio = ratio
	if isMove {
		p.Layout.Replace(depth, view)
	} else {
		p.Layout.Place(depth, view)
	}
}

func (p *SWFTreeProcessor) process(actions ActionList) (tag swftag.Tag, newActions ActionList) {
	tag = p.Current()
	if tag == nil {
		return nil, nil
	}

	switch node := tag.(type) {
	case *swftag.DefineMorphShape:
		if p.Loops > 0 {
			break
		}
		p.Objects.Add(MorphShapeDefinitionFromSWF(node.CharacterId, shapes.RectangleFromSWF(node.StartBounds), shapes.RectangleFromSWF(node.EndBounds), node.StartEdges.Records, node.EndEdges.Records, node.MorphFillStyles, node.MorphLineStyles))
	case *swftag.DefineMorphShape2:
		if p.Loops > 0 {
			break
		}
		p.Objects.Add(MorphShapeDefinitionFromSWF(node.CharacterId, shapes.RectangleFromSWF(node.StartBounds), shapes.RectangleFromSWF(node.EndBounds), node.StartEdges.Records, node.EndEdges.Records, node.MorphFillStyles, node.MorphLineStyles))
	case *swftag.DefineShape:
		p.Objects.Add(ShapeDefinitionFromSWF(node.ShapeId, shapes.RectangleFromSWF(node.ShapeBounds), node.Shapes.Records, node.Shapes.FillStyles, node.Shapes.LineStyles))
	case *swftag.DefineShape2:
		p.Objects.Add(ShapeDefinitionFromSWF(node.ShapeId, shapes.RectangleFromSWF(node.ShapeBounds), node.Shapes.Records, node.Shapes.FillStyles, node.Shapes.LineStyles))
	case *swftag.DefineShape3:
		p.Objects.Add(ShapeDefinitionFromSWF(node.ShapeId, shapes.RectangleFromSWF(node.ShapeBounds), node.Shapes.Records, node.Shapes.FillStyles, node.Shapes.LineStyles))
	case *swftag.DefineShape4:
		p.Objects.Add(ShapeDefinitionFromSWF(node.ShapeId, shapes.RectangleFromSWF(node.ShapeBounds), node.Shapes.Records, node.Shapes.FillStyles, node.Shapes.LineStyles))
	//TODO: case *swftag.DefineShape5:
	case *swftag.DefineSprite:
		if p.Loops > 0 {
			break
		}
		p.Objects.Add(&SpriteDefinition{
			ObjectId:  node.SpriteId,
			Processor: NewSWFTreeProcessor(node.SpriteId, node.ControlTags, p.Objects),
		})

	case *swftag.RemoveObject:
		//TODO: maybe replicate swftag.RemoveObject2 behavior?
		if o := p.Layout.Get(node.Depth); o != nil && o.GetObjectId() == node.CharacterId {
			p.Layout.Remove(node.Depth)
		} else {
			runtime.KeepAlive(o)
		}
	case *swftag.RemoveObject2:
		p.Layout.Remove(node.Depth)

	case *swftag.PlaceObject:
		var object ObjectDefinition
		if vl := p.Layout.Get(node.Depth); vl != nil {
			object = vl.Object
		}

		var transform *math.MatrixTransform
		if t := math.MatrixTransformFromSWF(node.Matrix); !t.IsIdentity() {
			transform = &t
		}

		var colorTransform *math.ColorTransform
		if node.Flag.HasColorTransform && node.ColorTransform != nil {
			t := math.ColorTransformFromSWF(*node.ColorTransform)
			colorTransform = &t
		}

		p.placeObject(object, node.Depth, 0, false, false, false, 0, transform, colorTransform)
	case *swftag.PlaceObject2:
		var object ObjectDefinition
		if node.Flag.HasCharacter {
			object = p.Objects[node.CharacterId]
		} else if vl := p.Layout.Get(node.Depth); vl != nil {
			object = vl.Object
		}

		var transform *math.MatrixTransform
		if node.Flag.HasMatrix {
			t := math.MatrixTransformFromSWF(node.Matrix)
			transform = &t
		}

		var colorTransform *math.ColorTransform
		if node.Flag.HasColorTransform {
			t := math.ColorTransformFromSWFAlpha(node.ColorTransform)
			colorTransform = &t
		}

		p.placeObject(object, node.Depth, node.ClipDepth, node.Flag.Move, node.Flag.HasRatio, node.Flag.HasClipDepth, float64(node.Ratio)/math2.MaxUint16, transform, colorTransform)
	case *swftag.PlaceObject3:
		//TODO: handle extra properties
		var object ObjectDefinition
		if node.Flag.HasCharacter {
			object = p.Objects[node.CharacterId]
		} else {
			object = p.Layout.Get(node.Depth).Object
		}

		var transform *math.MatrixTransform
		if node.Flag.HasMatrix {
			t := math.MatrixTransformFromSWF(node.Matrix)
			transform = &t
		}

		var colorTransform *math.ColorTransform
		if node.Flag.HasColorTransform {
			t := math.ColorTransformFromSWFAlpha(node.ColorTransform)
			colorTransform = &t
		}

		p.placeObject(object, node.Depth, node.ClipDepth, node.Flag.Move, node.Flag.HasRatio, node.Flag.HasClipDepth, float64(node.Ratio)/math2.MaxUint16, transform, colorTransform)

	case *swftag.ShowFrame:
	case *swftag.End:
	case *swftag.DoAction:
		for _, action := range node.Actions {
			switch action.ActionCode {
			case subtypes.ActionStop:
				actions = append(actions, &StopAction{})
			case subtypes.ActionPlay:
				actions = append(actions, &PlayAction{})
				//TODO ActionGotoFrame
			case subtypes.ActionNextFrame:
				actions = append(actions, &NextFrameAction{})
				//TODO ActionPreviousFrame

			}
		}
		//TODO DoInitAction

	}

	return tag, actions
}

func (p *SWFTreeProcessor) NextFrame() *ViewFrame {
	var actions ActionList
	if !p.Playing {
		return p.LastFrame
	}

	var node swftag.Tag
	for {
		node, actions = p.Process(actions)
		if node == nil {
			break
		}
		p.Next()

		if _, ok := node.(*swftag.ShowFrame); ok {
			break
		} else if _, ok := node.(*swftag.End); ok && p.Frame == 0 {
			break
		}
	}

	if node == nil { //Loop again
		p.Loops++
		p.Frame = 0
		p.Index = 0
		p.Layout = NewViewLayout(p.Layout.GetObjectId(), nil, nil)
		if p.LastFrame != nil {
			return p.NextFrame()
		}
		return nil
	}

	p.Frame++

	frame := p.Layout.NextFrame(actions)

	p.LastFrame = frame

	//TODO: this might need to be elsewhere?
	for _, action := range actions {
		switch action := action.(type) {
		case *StopAction:
			p.Playing = false
		case *PlayAction:
			p.Playing = true
		case *NextFrameAction:
			return p.NextFrame()
		default:
			_ = action

		}
	}

	return frame
}
