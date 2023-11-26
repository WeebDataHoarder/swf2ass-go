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

	Objects shapes.ObjectCollection

	Tags []swftag.Tag

	Index int

	Frame int64

	LastFrame *ViewFrame
	Playing   bool
	Loops     int

	processFunc func(actions ActionList) (tag swftag.Tag, newActions ActionList)
}

func NewSWFTreeProcessor(objectId uint16, tags []swftag.Tag, objects shapes.ObjectCollection) *SWFTreeProcessor {
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

func (p *SWFTreeProcessor) placeObject(object shapes.ObjectDefinition, depth, clipDepth uint16, isMove, hasRatio, hasClipDepth bool, ratio float64, transform *math.MatrixTransform, colorTransform *math.ColorTransform) {
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
			currentLayout.Properties.Ratio = ratio
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
	view.Properties.Ratio = ratio
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
		p.Objects.Add(MorphShapeDefinitionFromSWF(p.Objects, node.CharacterId, shapes.RectangleFromSWF(node.StartBounds), shapes.RectangleFromSWF(node.EndBounds), node.StartEdges.Records, node.EndEdges.Records, node.MorphFillStyles, node.MorphLineStyles))
	case *swftag.DefineMorphShape2:
		if p.Loops > 0 {
			break
		}
		p.Objects.Add(MorphShapeDefinitionFromSWF(p.Objects, node.CharacterId, shapes.RectangleFromSWF(node.StartBounds), shapes.RectangleFromSWF(node.EndBounds), node.StartEdges.Records, node.EndEdges.Records, node.MorphFillStyles, node.MorphLineStyles))
	case *swftag.DefineShape:
		if p.Loops > 0 {
			break
		}
		p.Objects.Add(ShapeDefinitionFromSWF(p.Objects, node.ShapeId, shapes.RectangleFromSWF(node.ShapeBounds), node.Shapes.Records, node.Shapes.FillStyles, node.Shapes.LineStyles))
	case *swftag.DefineShape2:
		if p.Loops > 0 {
			break
		}
		p.Objects.Add(ShapeDefinitionFromSWF(p.Objects, node.ShapeId, shapes.RectangleFromSWF(node.ShapeBounds), node.Shapes.Records, node.Shapes.FillStyles, node.Shapes.LineStyles))
	case *swftag.DefineShape3:
		if p.Loops > 0 {
			break
		}
		p.Objects.Add(ShapeDefinitionFromSWF(p.Objects, node.ShapeId, shapes.RectangleFromSWF(node.ShapeBounds), node.Shapes.Records, node.Shapes.FillStyles, node.Shapes.LineStyles))
	case *swftag.DefineShape4:
		if p.Loops > 0 {
			break
		}
		p.Objects.Add(ShapeDefinitionFromSWF(p.Objects, node.ShapeId, shapes.RectangleFromSWF(node.ShapeBounds), node.Shapes.Records, node.Shapes.FillStyles, node.Shapes.LineStyles))
	//TODO: case *swftag.DefineShape5:
	case *swftag.DefineSprite:
		if p.Loops > 0 {
			break
		}
		p.Objects.Add(&SpriteDefinition{
			ObjectId:  node.SpriteId,
			Processor: NewSWFTreeProcessor(node.SpriteId, node.ControlTags, p.Objects),
		})
	case *swftag.DefineText:
		ob := TextDefinitionFromSWF(p.Objects, node.CharacterId, node.Bounds, node.TextRecords, node.Matrix)
		if ob == nil {
			fmt.Printf("invalid text definition")
		} else {
			p.Objects.Add(ob)
		}
	case *swftag.DefineText2:
		ob := TextDefinitionFromSWF(p.Objects, node.CharacterId, node.Bounds, node.TextRecords, node.Matrix)
		if ob == nil {
			fmt.Printf("invalid text definition")
		} else {
			p.Objects.Add(ob)
		}
	case *swftag.DefineFont:
		if p.Loops > 0 {
			break
		}
		p.Objects.Add(FontDefinitionFromSWF(
			node.FontId,
			nil,
			false,
			false,
			false,
			false,
			nil,
			node.ShapeTable,
			node.OffsetTable,
			[]uint8{},
			node.Scale(),
		))
	case *swftag.DefineFontInfo:
		if p.Loops > 0 {
			break
		}

		ob := p.Objects.Get(node.FontId)
		if ob == nil {
			panic("font not found!")
		}
		if ob, ok := ob.(*FontDefinition); ok {
			ob.Name = node.FontName
			ob.Italic = node.Flag.Italic
			ob.Bold = node.Flag.Bold
			if node.Flag.WideCodes {
				if len(node.CodeTable16) != len(ob.Entries) {
					panic("wrong code count")
				}
				for i := range ob.Entries {
					ob.Entries[i].Code = uint32(node.CodeTable16[i])
				}
			} else {
				if len(node.CodeTable8) != len(ob.Entries) {
					panic("wrong code count")
				}
				for i := range ob.Entries {
					ob.Entries[i].Code = uint32(node.CodeTable8[i])
				}
			}
		} else {
			panic("object is not font definition!")
		}
	case *swftag.DefineFont2:
		if p.Loops > 0 {
			break
		}
		if node.Flag.WideOffsets {
			if node.Flag.WideCodes {
				p.Objects.Add(FontDefinitionFromSWF(
					node.FontId,
					node.FontName,
					node.Flag.HasLayout,
					true,
					node.Flag.Italic,
					node.Flag.Bold,
					node.FontBoundsTable,
					node.ShapeTable,
					node.OffsetTable32,
					node.CodeTable16,
					node.Scale(),
				))
			} else {
				p.Objects.Add(FontDefinitionFromSWF(
					node.FontId,
					node.FontName,
					node.Flag.HasLayout,
					true,
					node.Flag.Italic,
					node.Flag.Bold,
					node.FontBoundsTable,
					node.ShapeTable,
					node.OffsetTable32,
					node.CodeTable8,
					node.Scale(),
				))
			}
		} else {
			if node.Flag.WideCodes {
				p.Objects.Add(FontDefinitionFromSWF(
					node.FontId,
					node.FontName,
					node.Flag.HasLayout,
					true,
					node.Flag.Italic,
					node.Flag.Bold,
					node.FontBoundsTable,
					node.ShapeTable,
					node.OffsetTable16,
					node.CodeTable16,
					node.Scale(),
				))
			} else {
				p.Objects.Add(FontDefinitionFromSWF(
					node.FontId,
					node.FontName,
					node.Flag.HasLayout,
					true,
					node.Flag.Italic,
					node.Flag.Bold,
					node.FontBoundsTable,
					node.ShapeTable,
					node.OffsetTable16,
					node.CodeTable8,
					node.Scale(),
				))
			}
		}
	case *swftag.DefineFont3:
		if node.Flag.WideOffsets {
			if node.Flag.WideCodes {
				p.Objects.Add(FontDefinitionFromSWF(
					node.FontId,
					node.FontName,
					node.Flag.HasLayout,
					true,
					node.Flag.Italic,
					node.Flag.Bold,
					node.FontBoundsTable,
					node.ShapeTable,
					node.OffsetTable32,
					node.CodeTable16,
					node.Scale(),
				))
			} else {
				p.Objects.Add(FontDefinitionFromSWF(
					node.FontId,
					node.FontName,
					node.Flag.HasLayout,
					true,
					node.Flag.Italic,
					node.Flag.Bold,
					node.FontBoundsTable,
					node.ShapeTable,
					node.OffsetTable32,
					node.CodeTable8,
					node.Scale(),
				))
			}
		} else {
			if node.Flag.WideCodes {
				p.Objects.Add(FontDefinitionFromSWF(
					node.FontId,
					node.FontName,
					node.Flag.HasLayout,
					true,
					node.Flag.Italic,
					node.Flag.Bold,
					node.FontBoundsTable,
					node.ShapeTable,
					node.OffsetTable16,
					node.CodeTable16,
					node.Scale(),
				))
			} else {
				p.Objects.Add(FontDefinitionFromSWF(
					node.FontId,
					node.FontName,
					node.Flag.HasLayout,
					true,
					node.Flag.Italic,
					node.Flag.Bold,
					node.FontBoundsTable,
					node.ShapeTable,
					node.OffsetTable16,
					node.CodeTable8,
					node.Scale(),
				))
			}
		}
	case *swftag.DefineFont4:
		print(node)
	case *swftag.DefineBits:
		if p.Loops > 0 {
			break
		}
		fmt.Printf("Unsupported image: DefineBits\n")
	case *swftag.DefineBitsLossless:
		if p.Loops > 0 {
			break
		}
		bitDef := BitmapDefinitionFromSWFLossless(node.CharacterId, node.GetImage())
		if bitDef == nil {
			fmt.Printf("Unsupported lossless bitmap\n")
			break
		}
		p.Objects.Add(bitDef)
	case *swftag.DefineBitsLossless2:
		if p.Loops > 0 {
			break
		}
		bitDef := BitmapDefinitionFromSWFLossless(node.CharacterId, node.GetImage())
		if bitDef == nil {
			fmt.Printf("Unsupported lossless bitmap\n")
			break
		}
		p.Objects.Add(bitDef)
	case *swftag.DefineBitsJPEG2:
		if p.Loops > 0 {
			break
		}
		bitDef, err := BitmapDefinitionFromSWF(node.CharacterId, node.Data, nil)
		if err != nil {
			fmt.Printf("Unsupported bitmap: %s\n", err)
			break
		}
		p.Objects.Add(bitDef)
	case *swftag.DefineBitsJPEG3:
		if p.Loops > 0 {
			break
		}
		bitDef, err := BitmapDefinitionFromSWF(node.CharacterId, node.ImageData, node.GetAlphaData())
		if err != nil {
			fmt.Printf("Unsupported bitmap: %s\n", err)
			break
		}
		p.Objects.Add(bitDef)
	case *swftag.DefineBitsJPEG4:
		if p.Loops > 0 {
			break
		}
		bitDef, err := BitmapDefinitionFromSWF(node.CharacterId, node.ImageData, node.GetAlphaData())
		if err != nil {
			fmt.Printf("Unsupported bitmap: %s\n", err)
			break
		}
		p.Objects.Add(bitDef)
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
		var object shapes.ObjectDefinition
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
		var object shapes.ObjectDefinition
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
		var object shapes.ObjectDefinition
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

		if node.Flag.HasBlendMode {
			fmt.Printf("Unsupported blends!!!\n")
			switch node.BlendMode {
			case swftag.BlendOverlay:
				//fake it somewhat with half transparency for now, TODO: split underlying image in intersections and hardcode-apply this
				i := math.IdentityColorTransform()
				i.Multiply.Alpha = 128

				if colorTransform != nil {
					i = i.Combine(*colorTransform)
				}
				colorTransform = &i
			}
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
			default:
				//fmt.Printf("unhandled action %d\n", action.ActionCode)
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
