package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/records"
	stdmath "math"
)

type Ellipse[T ~float64 | ~int64] struct {
	Center, Radius math.Vector2[T]
}

func NewCircle[T ~float64 | ~int64](center math.Vector2[T], radius T) Ellipse[T] {
	return Ellipse[T]{
		Center: center,
		Radius: math.NewVector2(radius, radius),
	}
}

// ellipseC 0.5522847498307935
var ellipseC = (4 / float64(3)) * (stdmath.Sqrt(2) - 1)

func ellipseDrawQuarter(center, size math.Vector2[float64]) *records.CubicCurveRecord {
	return &records.CubicCurveRecord{
		Control1: math.NewVector2(center.X-size.X, center.Y-ellipseC*size.Y),
		Control2: math.NewVector2(center.X-ellipseC*size.X, center.Y-size.Y),
		Anchor:   math.NewVector2(center.X, center.Y-size.Y),
		Start:    math.NewVector2(center.X-size.X, center.Y),
	}
}

func (r Ellipse[T]) Draw() []records.Record {
	center := r.Center.Float64()
	radius := r.Radius.Float64()
	return []records.Record{
		ellipseDrawQuarter(center, math.NewVector2(-radius.X, radius.Y)),
		ellipseDrawQuarter(center, radius).Reverse(), //Reverse so paths connect
		ellipseDrawQuarter(center, math.NewVector2(radius.X, -radius.Y)),
		ellipseDrawQuarter(center, math.NewVector2(-radius.X, -radius.Y)).Reverse(),
	}
}
