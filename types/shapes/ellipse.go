package shapes

import (
	swftypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/records"
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

func ellipseDrawQuarter(center, size math.Vector2[float64]) *records.CubicCurveRecord {
	const c = 0.55228474983 // (4/3) * (sqrt(2) - 1)

	return &records.CubicCurveRecord{
		Control1: math.Vector2ToType[float64, swftypes.Twip](math.NewVector2(center.X-size.X, center.Y-c*size.Y).Multiply(swftypes.TwipFactor)),
		Control2: math.Vector2ToType[float64, swftypes.Twip](math.NewVector2(center.X-c*size.X, center.Y-size.Y).Multiply(swftypes.TwipFactor)),
		Anchor:   math.Vector2ToType[float64, swftypes.Twip](math.NewVector2(center.X, center.Y-size.Y).Multiply(swftypes.TwipFactor)),
		Start:    math.Vector2ToType[float64, swftypes.Twip](math.NewVector2(center.X-size.X, center.Y).Multiply(swftypes.TwipFactor)),
	}
}

func (r Ellipse[T]) Draw() []records.Record {
	var center, radius math.Vector2[float64]
	switch any(r.Center.X).(type) {
	case swftypes.Twip:
		center = math.Vector2ToType[T, float64](r.Center).Multiply(swftypes.TwipFactor)
		radius = math.Vector2ToType[T, float64](r.Radius).Multiply(swftypes.TwipFactor)
	case int64, float64:
		center = math.Vector2ToType[T, float64](r.Center)
		radius = math.Vector2ToType[T, float64](r.Radius)
	}
	return []records.Record{
		ellipseDrawQuarter(center, math.NewVector2(-radius.X, radius.Y)),
		ellipseDrawQuarter(center, radius).Reverse(), //Reverse so paths connect
		ellipseDrawQuarter(center, math.NewVector2(radius.X, -radius.Y)),
		ellipseDrawQuarter(center, math.NewVector2(-radius.X, -radius.Y)).Reverse(),
	}
}
