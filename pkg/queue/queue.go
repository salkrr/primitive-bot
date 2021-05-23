// Package queue implements a simple linked list based queue
package queue

import (
	"container/list"
	"sync"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
)

// Operation contains information needed to create
// primitive image.
type Operation struct {
	UserID  int64
	ImgPath string
	Config  primitive.Config
}

// Queue represents a linked list based queue.
// It contains Operation objects.
type Queue struct {
	elements *list.List
	mu       sync.Mutex
}

// New returns an initialized queue.
func New() *Queue {
	return &Queue{
		elements: list.New(),
		mu:       sync.Mutex{},
	}
}

// Enqueue adds element v to the end of the queue
// and returns its position.
func (q *Queue) Enqueue(v Operation) int {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.elements.PushBack(v)
	return q.elements.Len()
}

// Dequeue removes first element of the queue
// and returns it. If the queue is empty then
// second return parameter will be equal to false.
func (q *Queue) Dequeue() (Operation, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	e := q.elements.Front()
	if e == nil {
		return Operation{}, false
	}

	q.elements.Remove(e)
	return e.Value.(Operation), true
}

// Peek returns firsth element of the queue.
// If the queue is empty then second return parameter
// will be equal to false.
func (q *Queue) Peek() (Operation, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	e := q.elements.Front()
	if e == nil {
		return Operation{}, false
	}

	return e.Value.(Operation), true
}

// GetOperations returns operation with the given chatID and also the slice which
// contains positions of these operations.
func (q *Queue) GetOperations(userID int64) ([]Operation, []int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var operations []Operation
	var positions []int
	for i, e := 1, q.elements.Front(); e != nil; i, e = i+1, e.Next() {
		op := e.Value.(Operation)
		if op.UserID == userID {
			operations = append(operations, op)
			positions = append(positions, i)
		}
	}
	return operations, positions
}

// GetNumOperations returns the number of operations with specified chatID
// that are currently in the queue.
func (q *Queue) GetNumOperations(chatID int64) int {
	q.mu.Lock()
	defer q.mu.Unlock()

	counter := 0
	for i, e := 1, q.elements.Front(); e != nil; i, e = i+1, e.Next() {
		op := e.Value.(Operation)
		if op.UserID == chatID {
			counter++
		}
	}
	return counter
}
