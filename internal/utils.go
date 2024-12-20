package internal

import (
	"errors"
	"fmt"
)

func MapSlice[T any, V any](in []T, mapper func(T) V) []V {
	out := make([]V, len(in))
	for i, t := range in {
		out[i] = mapper(t)
	}
	return out
}

func MapSlicePE[T any, V any](in []T, mapper func(T) (*V, error)) ([]V, error) {
	out := make([]V, len(in))
	for i, t := range in {
		v, err := mapper(t)
		if err != nil {
			return nil, errors.Join(errors.New(fmt.Sprintf("on val %d", i)), err)
		}
		out[i] = *v
	}
	return out, nil
}
