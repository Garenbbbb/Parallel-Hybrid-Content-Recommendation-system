package deque

import (
	"Work-Stealing-Based-Parallel-Hybrid-Content-Recommendation-system/data"
	"sync"
)

type Task struct {
	ID    int
	Count int
}

type TaskCart struct {
	Task
	Info data.ShopingCart
}

type TaskItem struct {
	Task
	Info data.Content
}

type Deque[T any] struct {
	mu   sync.Mutex
	head *node[T]
	tail *node[T]
}

// node represents a node in the deque
type node[T any] struct {
	task T
	next *node[T]
	prev *node[T]
}

// PushFront adds a task to the front of the deque
func (d *Deque[T]) PushFront(task T) {
	d.mu.Lock()
	defer d.mu.Unlock()

	newNode := &node[T]{task: task}

	if d.head == nil {
		d.head = newNode
		d.tail = newNode
	} else {
		newNode.next = d.head
		d.head.prev = newNode
		d.head = newNode
	}
}

// PopFront removes and returns a task from the front of the deque
func (d *Deque[T]) PopFront() (T, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.head == nil {
		var empty T
		return empty, false
	}

	task := d.head.task
	d.head = d.head.next

	if d.head == nil {
		d.tail = nil
	} else {
		d.head.prev = nil
	}

	return task, true
}

// PushBack adds a task to the back of the deque
func (d *Deque[T]) PushBack(task T) {
	d.mu.Lock()
	defer d.mu.Unlock()

	newNode := &node[T]{task: task}

	if d.tail == nil {
		d.head = newNode
		d.tail = newNode
	} else {
		newNode.prev = d.tail
		d.tail.next = newNode
		d.tail = newNode
	}
}

// PopBack removes and returns a task from the back of the deque
func (d *Deque[T]) PopBack() (T, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.tail == nil {
		var empty T
		return empty, false
	}

	task := d.tail.task
	d.tail = d.tail.prev

	if d.tail == nil {
		d.head = nil
	} else {
		d.tail.next = nil
	}

	return task, true
}
