package util

type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
    if len(s.items) == 0 {
        panic("Stack underflow")
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}

func (s *Stack[T]) Peek(idx int) T {
    if idx >= len(s.items) {
        panic("peeking in wrong location of Stack")
    }
    return s.items[idx]
}
