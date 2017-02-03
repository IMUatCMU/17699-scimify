package adt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStack_Push(t *testing.T) {
	s := NewStack(5)
	assert.Equal(t, 0, s.Size())

	s.Push("1")
	assert.Equal(t, 1, s.Size())

	s.Push("2")
	assert.Equal(t, 2, s.Size())

	s.Push("3")
	assert.Equal(t, 3, s.Size())

	s.Push("4")
	assert.Equal(t, 4, s.Size())

	s.Push("5")
	assert.Equal(t, 5, s.Size())
}

func TestStack_Pop(t *testing.T) {
	s := NewStack(5)
	s.Push("1")
	s.Push("2")
	s.Push("3")
	s.Push("4")
	s.Push("5")

	assert.Equal(t, "5", s.Pop())
	assert.Equal(t, "4", s.Pop())
	assert.Equal(t, "3", s.Pop())
	assert.Equal(t, "2", s.Pop())
	assert.Equal(t, "1", s.Pop())
	assert.Nil(t, s.Pop())
}

func TestStack_Peek(t *testing.T) {
	s := NewStack(5)
	assert.Nil(t, s.Peek())

	s.Push("1")
	assert.Equal(t, "1", s.Peek())

	s.Push("2")
	assert.Equal(t, "2", s.Peek())

	s.Push("3")
	assert.Equal(t, "3", s.Peek())

	s.Push("4")
	assert.Equal(t, "4", s.Peek())

	s.Push("5")
	assert.Equal(t, "5", s.Peek())
}

func TestStack_Capacity(t *testing.T) {
	s := NewStack(5)
	assert.Equal(t, 5, s.Capacity())
}
