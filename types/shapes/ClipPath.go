package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"github.com/ctessum/geom"
	"slices"
)

type ClipPath struct {
	Clip ComplexPolygon
}

func NewClipPath(shape Shape) *ClipPath {
	return &ClipPath{
		Clip: ComplexPolygon{
			Pol: NewPolygonFromShape(shape),
		},
	}
}

func (c *ClipPath) AddShape(shape Shape) {
	c.Clip.Pol = c.Clip.Pol.Union(NewPolygonFromShape(shape))
}

func (c *ClipPath) GetShape() Shape {
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

// ClipShape Clips a shape, but attempts to recover original curved records
func (c *ClipPath) ClipShape(o Shape, recover bool) (r Shape) {
	if !recover {
		return c.Clip.Intersect(ComplexPolygon{
			Pol: NewPolygonFromShape(o),
		}).GetShape()
	}
	flatShape, correspondence := o.FlattenWithCorrespondence()
	outShape := c.Clip.Intersect(ComplexPolygon{
		Pol: NewPolygonFromShape(flatShape),
	}).GetShape()

	for i := 0; i < len(outShape); i++ {
		var found bool
		for j, e := range correspondence {
			if func() bool {
				k, l := i, 0
				for {
					if l >= len(e.Flattened) {
						return true
					}
					if k >= len(outShape) {
						return false
					}
					if !e.Flattened[l].Equals(outShape[k]) {
						return false
					}
					k++
					l++
				}
			}() {
				//They are the same! Append entry back
				i += len(e.Flattened) - 1
				r = append(r, e.Original)
				slices.Delete(correspondence, j, j+1)
				found = true
				break
			}
		}
		if !found {
			r = append(r, outShape[i])
		}
	}

	return r
}

func (c *ClipPath) Intersect(o *ClipPath) *ClipPath {
	return &ClipPath{
		Clip: c.Clip.Intersect(o.Clip),
	}
}
