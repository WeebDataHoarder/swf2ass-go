package shapes

import "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"

type StrokeConverter struct {
	HalfWidth  float64
	MiterLimit float64
	Cap        Capper
	Join       Joiner

	Position math.Vector2[float64]

	Segment [2]PathSegment[float64]
}

type Capper func(from, to math.Vector2[float64]) PathSegment[float64]

type Joiner func(from, to math.Vector2[float64]) PathSegment[float64]

var ButtCapper Capper = func(from, to math.Vector2[float64]) PathSegment[float64] {
	return PathSegment[float64]{
		{
			Pos:             from,
			IsBezierControl: false,
		},
		{
			Pos:             to,
			IsBezierControl: false,
		},
	}
}

var RoundCapper Capper = func(from, to math.Vector2[float64]) PathSegment[float64] {
	mid := from.AddVector(to).Divide(2)
	//TODO
	ellipseDrawQuarter(mid, from.SubVector(mid)).ToLineRecords(1)
	ellipseDrawQuarter(mid, to.SubVector(mid)).ToLineRecords(1)
	return PathSegment[float64]{
		{
			Pos:             from,
			IsBezierControl: false,
		},
		{
			Pos:             to,
			IsBezierControl: false,
		},
	}
}

var StraightJoiner Joiner = func(from, to math.Vector2[float64]) PathSegment[float64] {
	return PathSegment[float64]{
		{
			Pos:             from,
			IsBezierControl: false,
		},
		{
			Pos:             to,
			IsBezierControl: false,
		},
	}
}

func strokeMergeSegmentEnd(a, b PathSegment[float64]) PathSegment[float64] {
	if a.End().Equals(b.Start()) {
		return append(a, b[1:]...)
	} else if a.End().Equals(b.End()) {
		b.Flip()
		return append(a, b[1:]...)
	} else {
		panic("not joined!")
	}
}

func (c *StrokeConverter) checkInit(normal0, normal1 math.Vector2[float64]) {
	if len(c.Segment[0]) == 0 {
		//Init
		c.addPoint(c.Position, normal0, normal1)
	}
}

func (c *StrokeConverter) addPoint(v, normal0, normal1 math.Vector2[float64]) {

	start0 := c.Position.AddVector(normal0.Multiply(c.HalfWidth))
	start1 := c.Position.AddVector(normal1.Multiply(c.HalfWidth))
	point0 := v.AddVector(normal0.Multiply(c.HalfWidth))
	point1 := v.AddVector(normal1.Multiply(c.HalfWidth))

	if len(c.Segment[0]) > 0 {
		if !c.Segment[0].End().Equals(start0) {
			join0 := c.Join(c.Segment[0].End(), start0)
			c.Segment[0] = strokeMergeSegmentEnd(c.Segment[0], join0)
		}
		if !c.Segment[1].End().Equals(start1) {
			join0 := c.Join(c.Segment[1].End(), start1)
			c.Segment[1] = strokeMergeSegmentEnd(c.Segment[1], join0)
		}
	}

	c.Segment[0].AddPoint(VisitedPoint[float64]{
		Pos:             point0,
		IsBezierControl: false,
	})
	c.Segment[1].AddPoint(VisitedPoint[float64]{
		Pos:             point1,
		IsBezierControl: false,
	})
}

func (c *StrokeConverter) Line(v math.Vector2[float64]) {
	normal0, normal1 := v.SubVector(c.Position).Normals()
	normal0 = normal0.Normalize()
	normal1 = normal1.Normalize()
	c.checkInit(normal0, normal1)
	c.addPoint(v, normal0, normal1)
	c.Position = v
}

func (c *StrokeConverter) Close() PathSegment[float64] {
	if len(c.Segment[0]) <= 1 {
		return nil
	}
	topCap := c.Cap(c.Segment[0].Start(), c.Segment[1].Start())
	bottomCap := c.Cap(c.Segment[0].End(), c.Segment[1].End())
	segment := c.Segment[0]
	if !segment.TryMerge(&topCap, false) {
		panic("ouch")
	}
	if !segment.TryMerge(&bottomCap, false) {
		panic("ouch")
	}
	if !segment.TryMerge(&c.Segment[1], false) {
		panic("ouch")
	}

	c.Segment[0] = nil
	c.Segment[1] = nil
	return segment
}
