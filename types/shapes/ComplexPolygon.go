package shapes

import (
	swftypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/records"
	"github.com/ctessum/geom"
)

type ComplexPolygon struct {
	Pol geom.Polygonal
}

func (p ComplexPolygon) Intersect(o ComplexPolygon) ComplexPolygon {
	return ComplexPolygon{
		Pol: p.Pol.Intersection(o.Pol),
	}
}

func (p ComplexPolygon) GetShape() (r *Shape) {
	var edges []records.Record
	for _, pol := range p.Pol.Polygons() {
		for _, path := range pol {
			edges = append(edges, &records.LineRecord{
				To:    math.Vector2ToType[float64, swftypes.Twip](math.NewVector2(path[1].X, path[1].Y)),
				Start: math.Vector2ToType[float64, swftypes.Twip](math.NewVector2(path[0].X, path[0].Y)),
			})
			for _, point := range path[2:] {
				edges = append(edges, &records.LineRecord{
					To:    math.Vector2ToType[float64, swftypes.Twip](math.NewVector2(point.X, point.Y)),
					Start: edges[len(edges)-1].GetEnd(),
				})
			}
		}
	}
	return &Shape{
		Edges:  edges,
		IsFlat: true,
	}
}

func NewPolygonFromShape(shape *Shape) (g geom.Polygon) {
	flat := shape.Flatten()

	var edges []*records.LineRecord

	var lastEdge *records.LineRecord

	for _, record := range flat.Edges {
		if lastEdge != nil && !lastEdge.GetEnd().Equals(record.GetStart()) {
			g = append(g, NewPathFromEdges(edges))
			edges = edges[:0]
		}

		if lineRecord, ok := record.(*records.LineRecord); ok {
			edges = append(edges, lineRecord)
			lastEdge = lineRecord
		} else {
			panic("invalid record")
		}
	}

	if len(edges) > 0 {
		g = append(g, NewPathFromEdges(edges))
	}

	return g
}

func NewPathFromEdges(edges []*records.LineRecord) (p geom.Path) {
	p = make(geom.Path, 0, len(edges)+1)
	start := edges[0].Start.Float64()
	to := edges[0].To.Float64()
	p = append(p, geom.Point{
		X: start.X,
		Y: start.Y,
	}, geom.Point{
		X: to.X,
		Y: to.Y,
	})
	for _, e := range edges[1:] {
		to = e.To.Float64()
		p = append(p, geom.Point{
			X: to.X,
			Y: to.Y,
		})
	}

	//Close drawing if not closed
	if !start.Equals(to) {
		p = append(p, geom.Point{
			X: start.X,
			Y: start.Y,
		})
	}
	return p
}

func (p ComplexPolygon) Draw() []records.Record {
	return p.GetShape().Edges
}
