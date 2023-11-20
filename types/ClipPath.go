package types

import "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"

type ClipPath struct {
}

func NewClipPath(shape *Shape) *ClipPath {
	return &ClipPath{}
}

func (c *ClipPath) AddShape(shape *Shape) {

}

func (c *ClipPath) GetShape() *Shape {
	return &Shape{}
}

func (c *ClipPath) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) *ClipPath {
	return NewClipPath(c.GetShape().ApplyMatrixTransform(transform, applyTranslation))
}

func (c *ClipPath) Intersect(o *ClipPath) *ClipPath {
	//TODO: implement this
	return o
}
