package pkg

type Stack[T any] struct {
	Values []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{nil}
}

func (s *Stack[T]) Push(value T) {
	s.Values = append(s.Values, value)
}

func (s *Stack[T]) Pop() (T, bool) {
	var t, b = s.Peek()
	if b {
		s.Values = s.Values[:len(s.Values)-1]
	}
	return t, b
}

func (s *Stack[T]) Peek() (T, bool) {
	var t T
	if len(s.Values) == 0 {
		return t, false
	}
	return s.Values[len(s.Values)-1], true
}
