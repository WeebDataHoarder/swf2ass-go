package types

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	"github.com/ctessum/geom"
)

type ClipPath struct {
	Clip shapes.ComplexPolygon
}

func NewClipPath(shape *shapes.Shape) *ClipPath {
	if shape == nil {
		return &ClipPath{
			Clip: shapes.ComplexPolygon{
				Pol: shapes.NewPolygonFromShape(&shapes.Shape{}),
			},
		}
	}
	return &ClipPath{
		Clip: shapes.ComplexPolygon{
			Pol: shapes.NewPolygonFromShape(shape),
		},
	}
}

func (c *ClipPath) AddShape(shape *shapes.Shape) {
	c.Clip.Pol = c.Clip.Pol.Union(shapes.NewPolygonFromShape(shape))
}

func (c *ClipPath) GetShape() *shapes.Shape {
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
			Clip: shapes.ComplexPolygon{
				Pol: newPol,
			},
		}
	}
}

func (c *ClipPath) Intersect(o *ClipPath) *ClipPath {
	return &ClipPath{
		Clip: c.Clip.Intersect(o.Clip),
	}
}
