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

func (s *Stack[T]) Assign(r Stack[T]) {
    s.items = make([]T, len(r.items))
    for i, item := range r.items {
        s.items[i] = item
    }
}

func (s *Stack[T]) Copy() *Stack[T] {
    copyStack := new(Stack[T])
    i := 0
    for i < len(s.items) {
        copyStack.Push(s.items[i])
        i += 1
    }
    return copyStack
}
