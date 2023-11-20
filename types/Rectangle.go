package types

import "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"

type Rectangle[T ~float64 | ~int64] struct {
	TopLeft, BottomRight Vector2[T]
}

func (r Rectangle[T]) InBounds(pos Vector2[T]) bool {
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

func (r Rectangle[T]) Draw() []Record {
	var tl, br Vector2[types.Twip]
	switch any(r.TopLeft.X).(type) {
	case types.Twip:
		tl = Vector2ToType[T, types.Twip](r.TopLeft)
		br = Vector2ToType[T, types.Twip](r.BottomRight)
	case int64, float64:
		tl = Vector2ToType[T, types.Twip](r.TopLeft.Multiply(types.TwipFactor))
		br = Vector2ToType[T, types.Twip](r.BottomRight.Multiply(types.TwipFactor))
	}
	return []Record{
		&LineRecord{
			To:    NewVector2(tl.X, br.Y),
			Start: tl,
		},
		&LineRecord{
			To:    br,
			Start: NewVector2(tl.X, br.Y),
		},
		&LineRecord{
			To:    NewVector2(br.X, tl.Y),
			Start: br,
		},
		&LineRecord{
			To:    tl,
			Start: NewVector2(br.X, tl.Y),
		},
	}
}

func (r Rectangle[T]) DrawOpen() []Record {
	return r.Draw()[:3]
}
