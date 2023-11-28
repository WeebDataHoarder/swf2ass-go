package types

type Option[T any] struct {
	value *T
}

func (o Option[T]) Some() (T, bool) {
	if o.value != nil {
		return *o.value, true
	} else {
		var zero T
		return zero, false
	}
}

func (o Option[T]) Unwrap() T {
	if o.value == nil {
		panic("Option must have Some")
	}
	return *o.value
}

type Combinable[T any] interface {
	Combine(o T) T
}

// Combine Combines two Option or returns either if any is Option.Some but not the other
// f is not necessary if T implements Combinable
func (o Option[T]) Combine(other Option[T], f func(a, b T) Option[T]) (result Option[T]) {
	if a, ok := o.Some(); ok {
		if b, ok := other.Some(); ok {
			if c, ok := any(a).(Combinable[T]); ok {
				return Some(c.Combine(b))
			}
			return f(a, b)
		}
		return o
	} else if _, ok := other.Some(); ok {
		return other
	} else {
		return Option[T]{}
	}
}

func (o Option[T]) Pointer() *T {
	return o.value
}

func (o Option[T]) With(f func(T)) {
	v, ok := o.Some()
	if ok {
		f(v)
	}
}

func None[T any]() Option[T] {
	return Option[T]{}
}

func Some[T any](value T) Option[T] {
	return Option[T]{
		value: &value,
	}
}

func SomePointer[T any](value *T) Option[T] {
	return Option[T]{
		value: value,
	}
}

func SomeWith[T any](value T, ok bool) Option[T] {
	if !ok {
		return Option[T]{}
	}
	return Option[T]{
		value: &value,
	}
}

func SomeDefault[T any](o Option[T], def T) Option[T] {
	if _, ok := o.Some(); ok {
		return o
	}
	return Some(def)
}
