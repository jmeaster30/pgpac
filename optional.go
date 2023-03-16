package main

type Optional[T any] struct {
	hasValue bool
	stored   T
}

func None[T any]() Optional[T] {
	return Optional[T]{
		hasValue: false,
	}
}

func Some[T any](value T) Optional[T] {
	return Optional[T]{
		hasValue: true,
		stored:   value,
	}
}

func (o Optional[T]) Value() T {
	if o.hasValue {
		return o.stored
	}
	panic("Optional does not have a value :(")
}

func (o Optional[T]) ValueOrDefault(defaultValue T) T {
	if o.hasValue {
		return o.stored
	}
	return defaultValue
}

func (o Optional[T]) HasValue() bool {
	return o.hasValue
}
