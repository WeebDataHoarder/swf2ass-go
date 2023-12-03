package math

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf-go/types"
	"math"
	"reflect"
)

type Vector2[T ~int64 | ~float64] struct {
	X T
	Y T
}

func NewVector2[T ~int64 | ~float64](x, y T) Vector2[T] {
	return Vector2[T]{
		X: x,
		Y: y,
	}
}

var epsilon = math.Abs(float64(7.)/3 - float64(4.)/3 - float64(1.))

func (v Vector2[T]) Equals(b Vector2[T]) bool {
	switch any(v.X).(type) {
	case int64, types.Twip:
		return v.equalsInt(b)
	case float64:
		return v.equalsFloat(b)
	default:
		//slow path to find underlying type
		switch reflect.TypeOf(v.X).Kind() {
		case reflect.Int64:
			return v.equalsInt(b)
		case reflect.Float64:
			return v.equalsFloat(b)
		}
	}
	panic("unsupported type")
}

func (v Vector2[T]) equalsFloat(b Vector2[T]) bool {
	return v == b || ((math.Abs(float64(b.X-v.X))) <= epsilon && (math.Abs(float64(b.Y-v.Y))) <= epsilon)
}

func (v Vector2[T]) equalsInt(b Vector2[T]) bool {
	return v == b
}

func (v Vector2[T]) Distance(b Vector2[T]) float64 {
	return math.Sqrt(float64(v.SquaredDistance(b)))
}

func (v Vector2[T]) SquaredDistance(b Vector2[T]) T {
	r := v.SubVector(b)
	r = r.MultiplyVector(r)
	return r.X + r.Y
}

func (v Vector2[T]) MultiplyVector(b Vector2[T]) Vector2[T] {
	return Vector2[T]{
		X: v.X * b.X,
		Y: v.Y * b.Y,
	}
}

func (v Vector2[T]) DivideVector(b Vector2[T]) Vector2[T] {
	return Vector2[T]{
		X: v.X / b.X,
		Y: v.Y / b.Y,
	}
}

func (v Vector2[T]) Multiply(size T) Vector2[T] {
	return Vector2[T]{
		X: v.X * size,
		Y: v.Y * size,
	}
}

func (v Vector2[T]) Divide(size T) Vector2[T] {
	return Vector2[T]{
		X: v.X / size,
		Y: v.Y / size,
	}
}

func (v Vector2[T]) Invert() Vector2[T] {
	return Vector2[T]{
		X: v.Y,
		Y: v.X,
	}
}

func (v Vector2[T]) AddVector(b Vector2[T]) Vector2[T] {
	return Vector2[T]{
		X: v.X + b.X,
		Y: v.Y + b.Y,
	}
}

func (v Vector2[T]) SubVector(b Vector2[T]) Vector2[T] {
	return Vector2[T]{
		X: v.X - b.X,
		Y: v.Y - b.Y,
	}
}

func (v Vector2[T]) Normals() (a, b Vector2[T]) {
	return Vector2[T]{
			X: -v.Y,
			Y: v.X,
		}, Vector2[T]{
			X: v.Y,
			Y: -v.X,
		}
}

func (v Vector2[T]) Max(b Vector2[T]) Vector2[T] {
	return NewVector2(max(v.X, b.X), max(v.Y, b.Y))
}

func (v Vector2[T]) Min(b Vector2[T]) Vector2[T] {
	return NewVector2(min(v.X, b.X), min(v.Y, b.Y))
}

func (v Vector2[T]) Abs() Vector2[T] {
	x := v.X
	if x < 0 {
		x = -x
	}
	y := v.Y
	if y < 0 {
		y = -y
	}
	return Vector2[T]{
		X: x,
		Y: y,
	}
}

func (v Vector2[T]) Dot(b Vector2[T]) T {
	return v.X*b.X + v.Y*b.Y
}

func (v Vector2[T]) Length() float64 {
	return math.Sqrt(float64(v.SquaredLength()))
}

func (v Vector2[T]) SquaredLength() T {
	return v.X*v.X + v.Y*v.Y
}

func (v Vector2[T]) Normalize() Vector2[float64] {
	length := v.SquaredLength()
	if length > 0 {
		return v.Float64().Divide(math.Sqrt(float64(length)))
	}
	return Vector2[float64]{
		X: 0,
		Y: 0,
	}
}

func (v Vector2[T]) Float64() Vector2[float64] {
	if fX, ok := any(v.X).(float64er); ok {
		return Vector2[float64]{
			X: fX.Float64(),
			Y: any(v.Y).(float64er).Float64(),
		}
	}
	return Vector2ToType[T, float64](v)
}

func (v Vector2[T]) Int64() Vector2[int64] {
	return Vector2ToType[T, int64](v)
}

func Vector2ToType[T ~int64 | ~float64, T2 ~int64 | ~float64](v Vector2[T]) Vector2[T2] {
	var t T2
	switch any(t).(type) {
	case T: //same type
		return Vector2[T2]{
			X: T2(v.X),
			Y: T2(v.Y),
		}
	case fromFloat64er[T2]:
		return Vector2[T2]{
			//TODO: use unsafe?
			X: any(t).(fromFloat64er[T2]).FromFloat64(float64(v.X)),
			Y: any(t).(fromFloat64er[T2]).FromFloat64(float64(v.Y)),
		}
	case float64:
		if fX, ok := any(v.X).(float64er); ok {
			return Vector2[T2]{
				//TODO: use unsafe?
				X: T2(fX.Float64()),
				Y: T2(any(v.Y).(float64er).Float64()),
			}
		}
		return Vector2[T2]{
			X: T2(v.X),
			Y: T2(v.Y),
		}
	default:
		return Vector2[T2]{
			X: T2(v.X),
			Y: T2(v.Y),
		}
	}
}

type fromFloat64er[T ~int64 | ~float64] interface {
	FromFloat64(v float64) T
}

type float64er interface {
	Float64() float64
}

var stringerI = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

func (v Vector2[T]) String() string {
	typ := reflect.TypeOf(v.X)
	if typ.Implements(stringerI) {
		return fmt.Sprintf("Vector2[%s](%s, %s)", typ.Name(), any(v.X).(fmt.Stringer), any(v.Y).(fmt.Stringer))
	}
	switch typ.Kind() {
	case reflect.Int64:
		return fmt.Sprintf("Vector2[%s](%d, %d)", typ.Name(), int64(v.X), int64(v.Y))
	case reflect.Float64:
		return fmt.Sprintf("Vector2[%s](%f, %f)", typ.Name(), float64(v.X), float64(v.Y))
	}
	panic("unsupported type")
}
