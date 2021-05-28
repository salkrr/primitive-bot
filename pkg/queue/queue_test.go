package queue

import (
	"reflect"
	"testing"
)

func TestQueue_Enqueue(t *testing.T) {
	q := New()
	operations := []Operation{
		{UserID: 123456789},
		{UserID: 192837465},
		{UserID: 987654321},
	}

	for i, op := range operations {
		expectedPos := i + 1
		pos := q.Enqueue(op)
		if pos != expectedPos {
			t.Errorf("Got position %d; want %d", pos, expectedPos)
		}

		lastElem := q.elements.Back().Value.(Operation)
		if lastElem != op {
			t.Errorf("Last element in the queue: %+v; want %+v", lastElem, op)
		}
	}
}

func TestQueue_Dequeue(t *testing.T) {
	q := New()
	op := Operation{UserID: 123456789}
	q.Enqueue(op)

	v, ok := q.Dequeue()
	if !ok {
		t.Errorf("Queue is empty although it should contain one element.")
	}
	if v != op {
		t.Errorf("Got operation: %+v; want %+v", v, op)
	}

	_, ok = q.Dequeue()
	if ok {
		t.Errorf("Empty queue returned element.")
	}
}

func TestQueue_PeekReturnsFalseIfQueueIsEmpty(t *testing.T) {
	q := New()

	_, ok := q.Peek()
	if ok {
		t.Errorf("Empty queue returned element.")
	}
}

func TestQueue_PeekReturnsElementAndDoesNotDeleteItFromQueue(t *testing.T) {
	q := New()
	op := Operation{UserID: 123456789}
	q.Enqueue(op)

	v, ok := q.Peek()
	if !ok {
		t.Errorf("Empty queue returned element.")
	}
	if v != op {
		t.Errorf("Last element in the queue: %+v; want %+v", v, op)
	}

	if q.elements.Len() != 1 {
		t.Errorf("Number of elements in the queue changed.")
	}
}

func TestQueue_GetOperations(t *testing.T) {
	tests := []struct {
		name       string
		userID     int64
		queueState []Operation
		Operations []Operation
		Positions  []int
	}{
		{
			name:   "Correctly returns operations of the user",
			userID: 123456789,
			queueState: []Operation{
				{
					UserID:  123456789,
					ImgPath: "hello_world.jpg",
				},
				{
					UserID:  987654321,
					ImgPath: "brave_new_world.png",
				},
				{
					UserID:  123456789,
					ImgPath: "test.png",
				},
			},
			Operations: []Operation{
				{
					UserID:  123456789,
					ImgPath: "hello_world.jpg",
				},
				{
					UserID:  123456789,
					ImgPath: "test.png",
				},
			},
			Positions: []int{1, 3},
		},
		{
			name:   "Returns nothing if no operations of the user",
			userID: 123456789,
			queueState: []Operation{
				{
					UserID:  111111111,
					ImgPath: "hello_world.jpg",
				},
				{
					UserID:  987654321,
					ImgPath: "brave_new_world.png",
				},
				{
					UserID:  222222222,
					ImgPath: "test.png",
				},
			},
			Operations: nil,
			Positions:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New()
			for _, op := range tt.queueState {
				q.Enqueue(op)
			}

			operations, positions := q.GetOperations(tt.userID)
			if !reflect.DeepEqual(operations, tt.Operations) {
				t.Errorf("Got operations: %v; want %v", operations, tt.Operations)
			}
			if !reflect.DeepEqual(positions, tt.Positions) {
				t.Errorf("Got positions: %v; want %v", positions, tt.Positions)
			}
		})
	}
}

func TestQueue_GetNumOperations(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		queueState    []Operation
		numOperations int
	}{
		{
			name:   "Correctly returns number of user's operations",
			userID: 123456789,
			queueState: []Operation{
				{
					UserID:  123456789,
					ImgPath: "hello_world.jpg",
				},
				{
					UserID:  987654321,
					ImgPath: "brave_new_world.png",
				},
				{
					UserID:  123456789,
					ImgPath: "test.png",
				},
			},
			numOperations: 2,
		},
		{
			name:   "Returns zero if user doesn't have operations",
			userID: 123456789,
			queueState: []Operation{
				{
					UserID:  111111111,
					ImgPath: "hello_world.jpg",
				},
				{
					UserID:  987654321,
					ImgPath: "brave_new_world.png",
				},
				{
					UserID:  222222222,
					ImgPath: "test.png",
				},
			},
			numOperations: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New()
			for _, op := range tt.queueState {
				q.Enqueue(op)
			}

			n := q.GetNumOperations(tt.userID)
			if n != tt.numOperations {
				t.Errorf("Got %d; want %d", n, tt.numOperations)
			}
		})
	}
}
