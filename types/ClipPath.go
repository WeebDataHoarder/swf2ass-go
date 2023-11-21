package types

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type ClipPath struct {
}

func NewClipPath(shape *shapes.Shape) *ClipPath {
	return &ClipPath{}
}

func (c *ClipPath) AddShape(shape *shapes.Shape) {

}

func (c *ClipPath) GetShape() *shapes.Shape {
	return &shapes.Shape{}
}

func (c *ClipPath) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) *ClipPath {
	return NewClipPath(c.GetShape().ApplyMatrixTransform(transform, applyTranslation))
}

func (c *ClipPath) Intersect(o *ClipPath) *ClipPath {
	//TODO: implement this
	return o
}
