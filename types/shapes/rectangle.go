package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/records"
)

type Rectangle[T ~float64 | ~int64] struct {
	TopLeft, BottomRight math.Vector2[T]
}

func NewSquare[T ~float64 | ~int64](topLeft math.Vector2[T], size T) Rectangle[T] {
	return Rectangle[T]{
		TopLeft:     topLeft,
		BottomRight: topLeft.AddVector(math.NewVector2(size, size)),
	}
}

func (r Rectangle[T]) InBounds(pos math.Vector2[T]) bool {
	return pos.X >= r.TopLeft.X && pos.Y >= r.TopLeft.Y && pos.X <= r.BottomRight.X && pos.Y <= r.BottomRight.Y
}

func (r Rectangle[T]) Width() T {
	return r.BottomRight.X - r.TopLeft.X
}

func (r Rectangle[T]) Height() T {
	return r.BottomRight.Y - r.TopLeft.Y
}

func (r Rectangle[T]) Area() T {
	return r.Width() * r.Height()
}

func (r Rectangle[T]) Divide(size T) Rectangle[T] {
	return Rectangle[T]{
		TopLeft:     r.TopLeft.Divide(size),
		BottomRight: r.BottomRight.Divide(size),
	}
}

func (r Rectangle[T]) Multiply(size T) Rectangle[T] {
	return Rectangle[T]{
		TopLeft:     r.TopLeft.Multiply(size),
		BottomRight: r.BottomRight.Multiply(size),
	}
}

func (r Rectangle[T]) Draw() []records.Record {
	tl := math.Vector2ToType[T, types.Twip](r.TopLeft)
	br := math.Vector2ToType[T, types.Twip](r.BottomRight)
	return []records.Record{
		&records.LineRecord{
			To:    math.NewVector2(tl.X, br.Y),
			Start: tl,
		},
		&records.LineRecord{
			To:    br,
			Start: math.NewVector2(tl.X, br.Y),
		},
		&records.LineRecord{
			To:    math.NewVector2(br.X, tl.Y),
			Start: br,
		},
		&records.LineRecord{
			To:    tl,
			Start: math.NewVector2(br.X, tl.Y),
		},
	}
}

func RectangleToType[T ~int64 | ~float64, T2 ~int64 | ~float64](r Rectangle[T]) Rectangle[T2] {
	return Rectangle[T2]{
		TopLeft:     math.Vector2ToType[T, T2](r.TopLeft),
		BottomRight: math.Vector2ToType[T, T2](r.BottomRight),
	}
}

func (r Rectangle[T]) DrawOpen() []records.Record {
	return r.Draw()[:3]
}

func RectangleFromSWF(rect types.Rectangle) Rectangle[types.Twip] {
	return Rectangle[types.Twip]{
		TopLeft:     math.NewVector2(rect.Xmin, rect.Ymin),
		BottomRight: math.NewVector2(rect.Xmax, rect.Ymax),
	}
}
