package types

type ClipPath struct {
}

func NewClipPath(shape *Shape) *ClipPath {
	return &ClipPath{}
}

func (c *ClipPath) GetShape() *Shape {
	return nil
}

func (c *ClipPath) Intersect(o *ClipPath) *ClipPath {
	//TODO: implement this
	return o
}
