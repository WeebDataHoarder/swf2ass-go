package types

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	"golang.org/x/exp/maps"
	"slices"
)

type ViewFrame struct {
	ObjectId uint16

	DepthMap map[uint16]*ViewFrame

	DrawPathList Option[shapes.DrawPathList]

	ColorTransform  Option[math.ColorTransform]
	MatrixTransform Option[math.MatrixTransform]

	ClipDepth Option[uint16]
}

func NewClippingFrame(objectId, clipDepth uint16, list Option[shapes.DrawPathList]) *ViewFrame {
	return &ViewFrame{
		ObjectId:     objectId,
		ClipDepth:    Some(clipDepth),
		DrawPathList: list,
		DepthMap:     make(map[uint16]*ViewFrame),
	}
}

func NewViewFrame(objectId uint16, list Option[shapes.DrawPathList]) *ViewFrame {
	return &ViewFrame{
		ObjectId:     objectId,
		DrawPathList: list,
		DepthMap:     make(map[uint16]*ViewFrame),
	}
}

func (f *ViewFrame) AddChild(depth uint16, frame *ViewFrame) {
	if _, ok := f.DrawPathList.Some(); ok {
		panic("adding child to item with draw list")
	}
	f.DepthMap[depth] = frame
}

func (f *ViewFrame) Render(baseDepth uint16, depthChain Depth, parentColor Option[math.ColorTransform], parentMatrix Option[math.MatrixTransform]) RenderedFrame {
	depthChain = append(slices.Clone(depthChain), baseDepth)

	matrixTransform := parentMatrix.Combine(f.MatrixTransform, nil)

	colorTransform := parentColor.Combine(f.ColorTransform, nil)

	var renderedFrame RenderedFrame

	if dpl, ok := f.DrawPathList.Some(); ok {
		renderedFrame = append(renderedFrame, &RenderedObject{
			Depth:           depthChain,
			ObjectId:        f.ObjectId,
			DrawPathList:    dpl,
			Clip:            nil,
			ColorTransform:  SomeDefault(colorTransform, math.IdentityColorTransform()).Unwrap(),
			MatrixTransform: SomeDefault(matrixTransform, math.IdentityTransform()).Unwrap(),
		})
	} else {
		clipMap := make(map[uint16]*ViewFrame)
		clipPaths := make(map[uint16]*shapes.ClipPath)

		keys := maps.Keys(f.DepthMap)
		slices.Sort(keys)
		for _, depth := range keys {
			frame := f.DepthMap[depth]
			if _, isClipping := frame.ClipDepth.Some(); isClipping { //Process clips as they come
				clipMap[depth] = frame
				var clipPath *shapes.ClipPath
				for _, clipObject := range frame.Render(depth, depthChain, colorTransform, matrixTransform) {
					clipShape := shapes.NewClipPath(nil)
					for _, p := range clipObject.DrawPathList {
						if _, ok := p.Style.(*shapes.FillStyleRecord); ok { //Only clip with fills TODO: is this correct?
							clipShape.AddShape(p.Shape)
						}
					}

					if len(clipShape.GetShape()) > 0 {
						//translate into absolute coordinates
						clipShape = clipShape.ApplyMatrixTransform(clipObject.MatrixTransform, true)
						if clipPath == nil {
							clipPath = clipShape
						} else {
							clipPath = clipShape.Intersect(clipPath)
						}
					}
				}

				if clipPath != nil {
					clipPaths[depth] = clipPath
				} else {
					delete(clipMap, depth) //TODO: ????
				}
			}
		}

		clipMapKeys := maps.Keys(clipMap)
		slices.Sort(clipMapKeys)

		for _, depth := range keys {
			frame := f.DepthMap[depth]
			if _, isClipping := frame.ClipDepth.Some(); isClipping { //Already processed
				continue
			}
			var clipPath *shapes.ClipPath

			for _, clipDepth := range clipMapKeys {
				clip := clipMap[clipDepth]
				if clip.ClipDepth.Unwrap() > depth && clipDepth < depth {
					if clipPath == nil {
						clipPath = clipPaths[clipDepth]
					} else {
						clipPath = clipPaths[clipDepth].Intersect(clipPath)
					}
				}
			}

			for _, object := range frame.Render(depth, depthChain, colorTransform, matrixTransform) {
				if object.Clip != nil && clipPath != nil {
					object.Clip = object.Clip.Intersect(clipPath)
				} else if clipPath != nil {
					object.Clip = clipPath
				}

				renderedFrame = append(renderedFrame, object)
			}
		}
	}

	return renderedFrame
}
