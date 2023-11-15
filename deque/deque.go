package deque

import (
	"Work-Stealing-Based-Parallel-Hybrid-Content-Recommendation-system/data"
	"sync"
)

type Task struct {
	ID   int
	Task data.ShopingCart
}

type Deque struct {
	mu   sync.Mutex
	head *node
	tail *node
}

// node represents a node in the deque
type node struct {
	task Task
	next *node
	prev *node
}

// PushFront adds a task to the front of the deque
func (d *Deque) PushFront(task Task) {
	d.mu.Lock()
	defer d.mu.Unlock()

	newNode := &node{task: task}

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
func (d *Deque) PopFront() (Task, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.head == nil {
		return Task{}, false
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
func (d *Deque) PushBack(task Task) {
	d.mu.Lock()
	defer d.mu.Unlock()

	newNode := &node{task: task}

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
func (d *Deque) PopBack() (Task, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.tail == nil {
		return Task{}, false
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
