package adt

type Stack interface {
	Push(item interface{})
	Peek() interface{}
	Pop() interface{}
	Size() int
	Capacity() int
}

type stack struct {
	data []interface{}
}

func NewStack(cap int) *stack {
	return &stack{data: make([]interface{}, 0, cap)}
}

func (s *stack) Push(item interface{}) {
	s.data = append(s.data, item)
}

func (s *stack) Peek() interface{} {
	if s.Size() == 0 {
		return nil
	}
	return s.data[s.Size()-1]
}

func (s *stack) Pop() interface{} {
	if s.Size() == 0 {
		return nil
	}
	item := s.data[s.Size()-1]
	s.data = s.data[0 : s.Size()-1]
	return item
}

func (s *stack) Size() int {
	return len(s.data)
}

func (s *stack) Capacity() int {
	return cap(s.data)
}
