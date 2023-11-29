package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/records"
	"github.com/ctessum/polyclip-go"
	math2 "math"
	"slices"
)

type ComplexPolygon struct {
	Pol polyclip.Polygon
}

func (p ComplexPolygon) Merge(o ComplexPolygon) ComplexPolygon {
	return ComplexPolygon{
		Pol: p.Pol.Construct(polyclip.UNION, o.Pol),
	}
}

func (p ComplexPolygon) Intersect(o ComplexPolygon) ComplexPolygon {
	return ComplexPolygon{
		Pol: p.Pol.Construct(polyclip.INTERSECTION, o.Pol),
	}
}

const PolygonSimplifyTolerance = 0.01

func (p ComplexPolygon) GetShape() (r Shape) {
	err := fixOrientation(p.Pol)
	if err != nil {
		panic(err)
	}
	for _, contour := range p.Pol {
		r = append(r, records.LineRecord{
			To:    math.NewVector2(contour[1].X, contour[1].Y),
			Start: math.NewVector2(contour[0].X, contour[0].Y),
		})
		for _, point := range contour[2:] {
			r = append(r, records.LineRecord{
				To:    math.NewVector2(point.X, point.Y),
				Start: r[len(r)-1].GetEnd(),
			})
		}
	}
	return r
}

func NewPolygonFromShape(shape Shape) (g polyclip.Polygon) {
	flat := shape.Flatten()

	var edges []records.LineRecord

	var lastEdgePos *math.Vector2[float64]

	for _, record := range flat {
		if lastEdgePos != nil && !lastEdgePos.Equals(record.GetStart()) {
			g = append(g, NewContourFromEdges(edges))
			edges = edges[:0]
		}

		if lineRecord, ok := record.(records.LineRecord); ok {
			edges = append(edges, lineRecord)
			p := lineRecord.GetEnd()
			lastEdgePos = &p
		} else if moveRecord, ok := record.(records.MoveRecord); ok {
			g = append(g, NewContourFromEdges(edges))
			edges = edges[:0]
			p := moveRecord.GetEnd()
			lastEdgePos = &p
		} else {
			panic("invalid record")
		}
	}

	if len(edges) > 0 {
		g = append(g, NewContourFromEdges(edges))
	}

	return g
}

func NewContourFromEdges(edges []records.LineRecord) (p polyclip.Contour) {
	p = make(polyclip.Contour, 0, len(edges)+1)
	start := edges[0].Start
	to := edges[0].To

	p = append(p, polyclip.Point{
		X: start.X,
		Y: start.Y,
	})

	if !start.Equals(to) {
		p = append(p, polyclip.Point{
			X: to.X,
			Y: to.Y,
		})
	}

	for _, e := range edges[1:] {
		p = append(p, polyclip.Point{
			X: e.To.X,
			Y: e.To.Y,
		})
	}

	/*if p[0] == p[len(p)-1] { //closed
		p = p[:len(p)-1]
	}*/

	if p[0] != p[len(p)-1] { //not closed
		p = append(p, p[0])
	}

	return p
}

// isLeft: test if a point is Left|On|Right of an infinite 2D line.
//
//	Input:  three points P0, P1, and P2
//	Return: >0 for P2 left of the line through P0 to P1
//	      =0 for P2 on the line
//	      <0 for P2 right of the line
//	From http://geomalgorithms.com/a01-_area.html#isLeft()
func isLeft(P0, P1, P2 polyclip.Point) float64 {
	return (P1.X-P0.X)*(P2.Y-P0.Y) -
		(P2.X-P0.X)*(P1.Y-P0.Y)
}

// orientation: test the orientation of a simple 2D polygon
//
//	Input:  Point* V = an array of n+1 vertex points with V[n]=V[0]
//	Return: >0 for counterclockwise
//	        =0 for none (degenerate)
//	        <0 for clockwise
//	Note: this algorithm is faster than computing the signed area.
//	From http://geomalgorithms.com/a01-_area.html#orientation2D_Polygon()
func orientation(V polyclip.Polygon) []float64 {
	// first find rightmost lowest vertex of the polygon
	out := make([]float64, len(V))
	for j, r := range V {
		rmin := 0
		xmin := r[0].X
		ymin := r[0].Y
		for i, p := range r {
			if p.Y > ymin {
				continue
			} else if p.Y == ymin { // just as low
				if p.X < xmin { // and to left
					continue
				}
			}
			rmin = i // a new rightmost lowest vertex
			xmin = p.X
			ymin = p.Y
		}

		// test orientation at the rmin vertex
		// ccw <=> the edge leaving V[rmin] is left of the entering edge
		if rmin == 0 || rmin == len(r)-1 {
			out[j] = isLeft(r[len(r)-2], r[0], r[1])
		} else {
			out[j] = isLeft(r[rmin-1], r[rmin], r[rmin+1])
		}
	}
	return out
}

func polyInPoly(outer, inner polyclip.Contour) bool {
	for _, p := range inner {
		if pointInPoly(p, outer) == 0 {
			return false
		}
	}
	return true
}

const tolerance = 1.e-9

func floatEquals(f1, f2 float64) bool {
	//return (f1 == f2)
	return (f1 == f2) ||
		(math2.Abs(f1-f2)/math2.Abs(f1+f2) < tolerance)
}

// returns 0 if false, +1 if true, -1 if pt ON polygon boundary
// See "The Point in Polygon Problem for Arbitrary Polygons" by Hormann & Agathos
// http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.88.5498&rep=rep1&type=pdf
func pointInPoly(pt polyclip.Point, path polyclip.Contour) int {
	result := 0
	cnt := len(path)
	if cnt < 3 {
		return 0
	}
	ip := path[0]
	for i := 1; i <= cnt; i++ {
		var ipNext polyclip.Point
		if i == cnt {
			ipNext = path[0]
		} else {
			ipNext = path[i]
		}
		if floatEquals(ipNext.Y, pt.Y) {
			if floatEquals(ipNext.X, pt.X) || (floatEquals(ip.Y, pt.Y) &&
				((ipNext.X-pt.X > -tolerance) == (ip.X-pt.X < tolerance))) {
				return -1
			}
		}
		if (ip.Y-pt.Y < tolerance) != (ipNext.Y-pt.Y < tolerance) {
			if ip.X-pt.X >= -tolerance {
				if ipNext.X-pt.X > -tolerance {
					result = 1 - result
				} else {
					d := (ip.X-pt.X)*(ipNext.Y-pt.Y) -
						(ipNext.X-pt.X)*(ip.Y-pt.Y)
					if floatEquals(d, 0) {
						return -1
					} else if (d > -tolerance) == (ipNext.Y-ip.Y > -tolerance) {
						result = 1 - result
					}
				}
			} else {
				if ipNext.X-pt.X > -tolerance {
					d := (ip.X-pt.X)*(ipNext.Y-pt.Y) -
						(ipNext.X-pt.X)*(ip.Y-pt.Y)
					if floatEquals(d, 0) {
						return -1
					} else if (d > -tolerance) == (ipNext.Y-ip.Y > -tolerance) {
						result = 1 - result
					}
				}
			}
		}
		ip = ipNext
	}
	return result
}

func fixOrientation(p polyclip.Polygon) error {
	o := orientation(p)
	for i, inner := range p {
		numInside := 0
		for j, outer := range p {
			if i != j {
				if polyInPoly(outer, inner) {
					numInside++
				}
			}
		}
		if numInside%2 == 1 && o[i] > 0. {
			slices.Reverse(inner)
		} else if numInside%2 == 0 && o[i] < 0. {
			slices.Reverse(inner)
		}
	}
	return nil
}

func (p ComplexPolygon) Draw() Shape {
	return p.GetShape()
}
