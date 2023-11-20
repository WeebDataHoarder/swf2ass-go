package types

type RenderedObject struct {
	Depth    Depth
	ObjectId uint16

	DrawPathList DrawPathList

	Clip *ClipPath

	ColorTransform  ColorTransform
	MatrixTransform MatrixTransform
}

func (o *RenderedObject) GetDepth() Depth {
	if len(o.Depth) > 0 && o.Depth[0] == 0 {
		return o.Depth[1:]
	}
	return o.Depth
}
