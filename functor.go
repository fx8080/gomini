package gomini

import (
	"context"
)

type Functor[T any] struct {
	ctx   context.Context
	value T
	err   error
}

func (f *Functor[T]) Of(mod T, err error) *Functor[T] {
	return &Functor[T]{value: mod, err: err}
}
func (f *Functor[T]) Map(fn func(T) (T, error)) *Functor[T] {
	if f.err != nil {
		return f
	}
	return f.Of(fn(f.value))
}
func (f *Functor[T]) Join() (T, error) {
	return f.value, f.err
}
