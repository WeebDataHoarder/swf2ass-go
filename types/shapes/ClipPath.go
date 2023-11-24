package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"github.com/ctessum/geom"
)

type ClipPath struct {
	Clip ComplexPolygon
}

func NewClipPath(shape *Shape) *ClipPath {
	if shape == nil {
		return &ClipPath{
			Clip: ComplexPolygon{
				Pol: NewPolygonFromShape(&Shape{}),
			},
		}
	}
	return &ClipPath{
		Clip: ComplexPolygon{
			Pol: NewPolygonFromShape(shape),
		},
	}
}

func (c *ClipPath) AddShape(shape *Shape) {
	c.Clip.Pol = c.Clip.Pol.Union(NewPolygonFromShape(shape))
}

func (c *ClipPath) GetShape() *Shape {
	return c.Clip.GetShape()
}

func (c *ClipPath) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) *ClipPath {
	pol, err := c.Clip.Pol.Transform(func(X, Y float64) (x, y float64, err error) {
		out := transform.ApplyToVector(math.NewVector2(X, Y), applyTranslation)
		return out.X, out.Y, nil
	})
	if err != nil {
		panic(err)
	}
	if newPol, ok := pol.(geom.Polygonal); !ok {
		panic("invalid result")
	} else {
		return &ClipPath{
			Clip: ComplexPolygon{
				Pol: newPol,
			},
		}
	}
}

func (c *ClipPath) Merge(o *ClipPath) *ClipPath {
	return &ClipPath{
		Clip: c.Clip.Merge(o.Clip),
	}
}

func (c *ClipPath) ClipShape(o *Shape) *Shape {
	return c.Clip.Intersect(ComplexPolygon{
		Pol: NewPolygonFromShape(o),
	}).GetShape()
}

func (c *ClipPath) Intersect(o *ClipPath) *ClipPath {
	return &ClipPath{
		Clip: c.Clip.Intersect(o.Clip),
	}
}
