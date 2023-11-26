package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/records"
	"github.com/ctessum/geom"
)

type ComplexPolygon struct {
	Pol geom.Polygonal
}

func (p ComplexPolygon) Merge(o ComplexPolygon) ComplexPolygon {
	return ComplexPolygon{
		Pol: p.Pol.Union(o.Pol),
	}
}

func (p ComplexPolygon) Intersect(o ComplexPolygon) ComplexPolygon {
	return ComplexPolygon{
		Pol: p.Pol.Intersection(o.Pol),
	}
}

const PolygonSimplifyTolerance = 0.01

func (p ComplexPolygon) GetShape() (r Shape) {
	for _, pol := range p.Pol.Polygons() {
		for _, path := range pol.Simplify(PolygonSimplifyTolerance).(geom.Polygon) {
			//pol = pol.Simplify(PolygonSimplifyTolerance).(geom.Polygon)
			r = append(r, records.LineRecord{
				To:    math.NewVector2(path[1].X, path[1].Y),
				Start: math.NewVector2(path[0].X, path[0].Y),
			})
			for _, point := range path[2:] {
				r = append(r, records.LineRecord{
					To:    math.NewVector2(point.X, point.Y),
					Start: r[len(r)-1].GetEnd(),
				})
			}
		}
	}
	return r
}

func NewPolygonFromShape(shape Shape) (g geom.Polygon) {
	flat := shape.Flatten()

	var edges []records.LineRecord

	var lastEdge *records.LineRecord

	for _, record := range flat {
		if lastEdge != nil && !lastEdge.GetEnd().Equals(record.GetStart()) {
			g = append(g, NewPathFromEdges(edges))
			edges = edges[:0]
		}

		if lineRecord, ok := record.(records.LineRecord); ok {
			edges = append(edges, lineRecord)
			lastEdge = &lineRecord
		} else {
			panic("invalid record")
		}
	}

	if len(edges) > 0 {
		g = append(g, NewPathFromEdges(edges))
	}

	return g
}

func NewPathFromEdges(edges []records.LineRecord) (p geom.Path) {
	p = make(geom.Path, 0, len(edges)+1)
	start := edges[0].Start
	to := edges[0].To

	p = append(p, geom.Point{
		X: start.X,
		Y: start.Y,
	})

	if !start.Equals(to) {
		p = append(p, geom.Point{
			X: to.X,
			Y: to.Y,
		})
	}

	for _, e := range edges[1:] {
		p = append(p, geom.Point{
			X: e.To.X,
			Y: e.To.Y,
		})
	}

	return p
}

func (p ComplexPolygon) Draw() Shape {
	return p.GetShape()
}
