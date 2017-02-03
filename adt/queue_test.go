package adt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueue_Capacity(t *testing.T) {
	q := NewQueue(10)
	assert.Equal(t, 10, q.Capacity())
}

func TestQueue_Offer(t *testing.T) {
	q := NewQueue(5)
	assert.Equal(t, 0, q.Size())

	q.Offer("1")
	assert.Equal(t, 1, q.Size())

	q.Offer("2")
	assert.Equal(t, 2, q.Size())

	q.Offer("3")
	assert.Equal(t, 3, q.Size())

	q.Offer("4")
	assert.Equal(t, 4, q.Size())

	q.Offer("5")
	assert.Equal(t, 5, q.Size())
}

func TestQueue_Peek(t *testing.T) {
	q := NewQueue(5)
	assert.Nil(t, q.Peek())

	q.Offer("1")
	assert.Equal(t, "1", q.Peek())

	q.Offer("2")
	assert.Equal(t, "1", q.Peek())

	q.Offer("3")
	assert.Equal(t, "1", q.Peek())

	q.Offer("4")
	assert.Equal(t, "1", q.Peek())

	q.Offer("5")
	assert.Equal(t, "1", q.Peek())
}

func TestQueue_Poll(t *testing.T) {
	q := NewQueue(5)
	q.Offer("1")
	q.Offer("2")
	q.Offer("3")
	q.Offer("4")
	q.Offer("5")

	assert.Equal(t, "1", q.Poll())
	assert.Equal(t, "2", q.Poll())
	assert.Equal(t, "3", q.Poll())
	assert.Equal(t, "4", q.Poll())
	assert.Equal(t, "5", q.Poll())
	assert.Nil(t, q.Poll())
}
