package reader

import (
	"dataProcessor/pkg/models"
	"sync"
)

type node struct {
	value models.ServiceStatus
	next  *node
}

type ConcurrentQueue struct {
	head *node
	tail *node
	mu   sync.Mutex
}

func NewConcurrentQueue() *ConcurrentQueue {
	return &ConcurrentQueue{
		head: nil,
		tail: nil,
		mu:   sync.Mutex{},
	}
}

func (q *ConcurrentQueue) Dequeue() *models.ServiceStatus {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.head == nil {
		return nil
	}

	removed := q.head.value

	if q.head == q.tail {
		q.head = nil
		q.tail = nil
	} else {
		q.head = q.head.next
	}

	return &removed
}

func (q *ConcurrentQueue) Enqueue(msg *models.ServiceStatus) {
	q.mu.Lock()
	defer q.mu.Unlock()

	newNode := &node{
		value: *msg,
		next:  nil,
	}

	if q.tail == nil {
		q.head = newNode
		q.tail = newNode
	} else {
		q.tail.next = newNode
		q.tail = newNode
	}
}
