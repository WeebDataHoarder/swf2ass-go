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

	DrawPathList *shapes.DrawPathList

	ColorTransform  *math.ColorTransform
	MatrixTransform *math.MatrixTransform

	IsClipping bool
	ClipDepth  uint16
}

func NewClippingFrame(objectId, clipDepth uint16, list *shapes.DrawPathList) *ViewFrame {
	return &ViewFrame{
		ObjectId:     objectId,
		ClipDepth:    clipDepth,
		DrawPathList: list,
		DepthMap:     make(map[uint16]*ViewFrame),
		IsClipping:   true,
	}
}

func NewViewFrame(objectId uint16, list *shapes.DrawPathList) *ViewFrame {
	return &ViewFrame{
		ObjectId:     objectId,
		DrawPathList: list,
		DepthMap:     make(map[uint16]*ViewFrame),
	}
}

func (f *ViewFrame) AddChild(depth uint16, frame *ViewFrame) {
	if f.DrawPathList != nil {
		panic("adding child to item with draw list")
	}
	f.DepthMap[depth] = frame
}

func (f *ViewFrame) Render(baseDepth uint16, depthChain Depth, parentColor *math.ColorTransform, parentMatrix *math.MatrixTransform) RenderedFrame {
	depthChain = slices.Clone(depthChain)
	depthChain = append(depthChain, baseDepth)

	matrixTransform := math.IdentityTransform()
	if f.MatrixTransform != nil {
		if parentMatrix != nil {
			matrixTransform = parentMatrix.Multiply(*f.MatrixTransform)
		} else {
			matrixTransform = *f.MatrixTransform
		}
	} else if parentMatrix != nil {
		matrixTransform = *parentMatrix
	}

	colorTransform := math.IdentityColorTransform()
	if f.ColorTransform != nil {
		if parentColor != nil {
			colorTransform = parentColor.Combine(*f.ColorTransform)
		} else {
			colorTransform = *f.ColorTransform
		}
	} else if parentColor != nil {
		colorTransform = *parentColor
	}

	var renderedFrame RenderedFrame

	if f.DrawPathList != nil {
		renderedFrame = append(renderedFrame, &RenderedObject{
			Depth:           depthChain,
			ObjectId:        f.ObjectId,
			DrawPathList:    *f.DrawPathList,
			Clip:            nil,
			ColorTransform:  colorTransform,
			MatrixTransform: matrixTransform,
		})
	} else {
		clipMap := make(map[uint16]*ViewFrame)
		clipPaths := make(map[uint16]*shapes.ClipPath)

		matrixTransform := &matrixTransform
		if matrixTransform.IsIdentity() {
			matrixTransform = nil
		}
		colorTransform := &colorTransform
		if colorTransform.IsIdentity() {
			colorTransform = nil
		}

		keys := maps.Keys(f.DepthMap)
		slices.Sort(keys)
		for _, depth := range keys {
			frame := f.DepthMap[depth]
			if frame.IsClipping { //Process clips as they come
				clipMap[depth] = frame
				var clipPath *shapes.ClipPath
				for _, clipObject := range frame.Render(depth, depthChain, colorTransform, matrixTransform) {
					clipShape := shapes.NewClipPath(nil)
					for _, p := range clipObject.DrawPathList {
						if _, ok := p.Style.(*shapes.FillStyleRecord); ok { //Only clip with fills TODO: is this correct?
							if p.Clip != nil {
								clipShape.AddShape(p.Clip.ClipShape(p.Commands, false))
							} else {
								clipShape.AddShape(p.Commands)
							}
						}
					}

					if len(clipShape.GetShape()) > 0 {
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
			if frame.IsClipping { //Already processed
				continue
			}
			var clipPath *shapes.ClipPath

			for _, clipDepth := range clipMapKeys {
				clip := clipMap[clipDepth]
				if clip.ClipDepth > depth && clipDepth < depth {
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
