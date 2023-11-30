package types

import (
	"fmt"
	swftag "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	math2 "math"
	"runtime"
	"slices"
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

	Version uint8

	JPEGTables []byte

	processFunc func(actions ActionList) (tag swftag.Tag, newActions ActionList)
}

func NewSWFTreeProcessor(objectId uint16, tags []swftag.Tag, objects shapes.ObjectCollection, version uint8) *SWFTreeProcessor {
	return &SWFTreeProcessor{
		Objects: objects,
		Frame:   0,
		Tags:    tags,
		Layout:  NewViewLayout(objectId, nil, nil),
		Playing: true,
		Version: version,
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

type PlaceAction uint8

func PlaceActionFromFlags(hasCharacter, isMove bool) PlaceAction {
	var action PlaceAction
	if hasCharacter && !isMove {
		action = ActionPlace
	} else if hasCharacter && isMove {
		action = ActionReplace
	} else if !hasCharacter && isMove {
		action = ActionModify
	} else {
		panic("invalid action")
	}
	return action
}

const (
	ActionPlace = PlaceAction(iota)
	ActionReplace
	ActionModify
)

type placeObjectData struct {
	Action         PlaceAction
	Depth          uint16
	ClipDepth      Option[uint16]
	Ratio          Option[float64]
	Transform      Option[math.MatrixTransform]
	ColorTransform Option[math.ColorTransform]
	Visible        Option[bool]
	BlendMode      Option[swftag.BlendMode]
}

func (p *SWFTreeProcessor) applyPlaceObject(layout *ViewLayout, data placeObjectData) {
	data.Transform.With(func(transform math.MatrixTransform) {
		layout.MatrixTransform = Some(transform)
	})
	data.ColorTransform.With(func(colorTransform math.ColorTransform) {
		layout.ColorTransform = Some(colorTransform)
	})
	data.Ratio.With(func(ratio float64) {
		layout.Properties.Ratio = ratio
	})

	if p.Version >= 11 {
		data.Visible.With(func(b bool) {
			layout.Properties.Visible = b
		})
		//todo: background color
	}
	//todo: filters
}

func (p *SWFTreeProcessor) placeObject(object shapes.ObjectDefinition, data placeObjectData) {
	if object == nil {
		//TODO: place bogus element
		fmt.Printf("Object at depth:%d not found\n", data.Depth)
		p.Layout.Remove(data.Depth)
		return
	}

	data.BlendMode.With(func(mode swftag.BlendMode) {
		fmt.Printf("Unsupported blends!!!\n")
		switch mode {
		case swftag.BlendOverlay:
			//fake it somewhat with half transparency for now, TODO: split underlying image in intersections and hardcode-apply this
			i := math.IdentityColorTransform()
			i.Multiply.Alpha = 128

			data.ColorTransform.With(func(transform math.ColorTransform) {
				i = i.Combine(transform)
			})

			data.ColorTransform = Some(i)
		}
	})

	switch data.Action {
	case ActionPlace:
		var view *ViewLayout
		if clipDepth, ok := data.ClipDepth.Some(); ok {
			view = NewClippingViewLayout(object.GetObjectId(), clipDepth, object.GetSafeObject(), p.Layout)
		} else {
			view = NewViewLayout(object.GetObjectId(), object.GetSafeObject(), p.Layout)
		}
		view.Properties.PlaceFrame = p.Frame
		p.applyPlaceObject(view, data)
		p.Layout.Place(data.Depth, view)
	case ActionReplace:
		var view *ViewLayout
		if clipDepth, ok := data.ClipDepth.Some(); ok {
			view = NewClippingViewLayout(object.GetObjectId(), clipDepth, object.GetSafeObject(), p.Layout)
		} else {
			view = NewViewLayout(object.GetObjectId(), object.GetSafeObject(), p.Layout)
		}
		view.Properties.PlaceFrame = p.Frame
		p.applyPlaceObject(view, data)
		p.Layout.Replace(data.Depth, view)
	case ActionModify:
		if currentLayout := p.Layout.Get(data.Depth); currentLayout != nil && currentLayout.GetObjectId() == object.GetObjectId() {
			p.applyPlaceObject(currentLayout, data)
		}
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
		p.Objects.Add(SpriteDefinitionFromSWF(node.SpriteId, int(node.FrameCount), NewSWFTreeProcessor(node.SpriteId, node.ControlTags, p.Objects, p.Version)))
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
	case *swftag.DefineBitsLossless:
		if p.Loops > 0 {
			break
		}
		im, err := node.GetImage()
		if err != nil {
			fmt.Printf("Unsupported lossless bitmap: %s\n", err)
			break
		}
		bitDef := BitmapDefinitionFromSWFLossless(node.CharacterId, im)
		if bitDef == nil {
			fmt.Printf("Unsupported lossless bitmap\n")
			break
		}
		p.Objects.Add(bitDef)
	case *swftag.DefineBitsLossless2:
		if p.Loops > 0 {
			break
		}
		im, err := node.GetImage()
		if err != nil {
			fmt.Printf("Unsupported lossless bitmap: %s\n", err)
			break
		}
		bitDef := BitmapDefinitionFromSWFLossless(node.CharacterId, im)
		if bitDef == nil {
			fmt.Printf("Unsupported lossless bitmap\n")
			break
		}
		p.Objects.Add(bitDef)
	case *swftag.DefineBits:
		if p.Loops > 0 {
			break
		}
		if p == nil {
			panic("todo: DefineBits within sprite??")
		}
		data := slices.Clone(p.JPEGTables)
		data = append(data, node.Data...)
		bitDef, err := BitmapDefinitionFromSWF(node.CharacterId, data, nil)
		if err != nil {
			fmt.Printf("Unsupported bitmap: %s\n", err)
			break
		}
		p.Objects.Add(bitDef)
	case *swftag.JPEGTables:
		if p.Loops > 0 {
			break
		}
		p.JPEGTables = node.Data
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
		object := p.Objects.Get(node.CharacterId)

		var colorTransform math.ColorTransform
		if node.Flag.HasColorTransform && node.ColorTransform != nil {
			t := math.ColorTransformFromSWF(*node.ColorTransform)
			colorTransform = t
		}

		transform := math.MatrixTransformFromSWF(node.Matrix, 1)

		p.placeObject(object, placeObjectData{
			Action:         ActionPlace,
			Depth:          node.Depth,
			ClipDepth:      None[uint16](),
			Ratio:          None[float64](),
			Transform:      SomeWith(transform, !transform.IsIdentity()),
			ColorTransform: SomeWith(colorTransform, node.Flag.HasColorTransform),
			Visible:        None[bool](),
		})
	case *swftag.PlaceObject2:
		var object shapes.ObjectDefinition
		if node.Flag.HasCharacter {
			object = p.Objects.Get(node.CharacterId)
		} else if vl := p.Layout.Get(node.Depth); vl != nil {
			object = vl.Object
		}

		p.placeObject(object, placeObjectData{
			Action:         PlaceActionFromFlags(node.Flag.HasCharacter, node.Flag.Move),
			Depth:          node.Depth,
			ClipDepth:      SomeWith(node.ClipDepth, node.Flag.HasClipDepth),
			Ratio:          SomeWith(float64(node.Ratio)/math2.MaxUint16, node.Flag.HasRatio),
			Transform:      SomeWith(math.MatrixTransformFromSWF(node.Matrix, 1), node.Flag.HasMatrix),
			ColorTransform: SomeWith(math.ColorTransformFromSWFAlpha(node.ColorTransform), node.Flag.HasColorTransform),
			Visible:        None[bool](),
		})
	case *swftag.PlaceObject3:
		//TODO: handle extra properties
		var object shapes.ObjectDefinition
		if node.Flag.HasCharacter {
			object = p.Objects.Get(node.CharacterId)
		} else {
			object = p.Layout.Get(node.Depth).Object
		}

		p.placeObject(object, placeObjectData{
			Action:         PlaceActionFromFlags(node.Flag.HasCharacter, node.Flag.Move),
			Depth:          node.Depth,
			ClipDepth:      SomeWith(node.ClipDepth, node.Flag.HasClipDepth),
			Ratio:          SomeWith(float64(node.Ratio)/math2.MaxUint16, node.Flag.HasRatio),
			Transform:      SomeWith(math.MatrixTransformFromSWF(node.Matrix, 1), node.Flag.HasMatrix),
			ColorTransform: SomeWith(math.ColorTransformFromSWFAlpha(node.ColorTransform), node.Flag.HasColorTransform),
			Visible:        SomeWith(node.Visible > 0, node.Flag.HasVisible),
			BlendMode:      SomeWith(node.BlendMode, node.Flag.HasBlendMode),
		})
	case *swftag.PlaceObject4:
		//TODO: handle extra properties
		var object shapes.ObjectDefinition
		if node.Flag.HasCharacter {
			object = p.Objects.Get(node.CharacterId)
		} else {
			object = p.Layout.Get(node.Depth).Object
		}

		p.placeObject(object, placeObjectData{
			Action:         PlaceActionFromFlags(node.Flag.HasCharacter, node.Flag.Move),
			Depth:          node.Depth,
			ClipDepth:      SomeWith(node.ClipDepth, node.Flag.HasClipDepth),
			Ratio:          SomeWith(float64(node.Ratio)/math2.MaxUint16, node.Flag.HasRatio),
			Transform:      SomeWith(math.MatrixTransformFromSWF(node.Matrix, 1), node.Flag.HasMatrix),
			ColorTransform: SomeWith(math.ColorTransformFromSWFAlpha(node.ColorTransform), node.Flag.HasColorTransform),
			Visible:        SomeWith(node.Visible > 0, node.Flag.HasVisible),
			BlendMode:      SomeWith(node.BlendMode, node.Flag.HasBlendMode),
		})
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

func (p *SWFTreeProcessor) NextFrame(loop bool) *ViewFrame {
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

	if node == nil {
		// We are done looping, check if we still need to keep playback
		p.Loops++
		//p.Frame = 0
		p.Index = 0
		//p.Layout = NewViewLayout(p.Layout.GetObjectId(), nil, nil)
		if loop && p.LastFrame != nil {
			return p.NextFrame(loop)
		}
		return nil
	}

	frame := p.Layout.NextFrame(p.Frame, actions)

	p.Frame++

	p.LastFrame = frame

	//TODO: this might need to be elsewhere?
	for _, action := range actions {
		switch action := action.(type) {
		case *StopAction:
			p.Playing = false
		case *PlayAction:
			p.Playing = true
		case *NextFrameAction:
			return p.NextFrame(loop)
		default:
			_ = action

		}
	}

	return frame
}
