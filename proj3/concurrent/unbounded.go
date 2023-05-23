package concurrent

import (
	"sync"
)

/**** YOU CANNOT MODIFY ANY OF THE FOLLOWING INTERFACES/TYPES ********/
type Task interface{}

type DEQueue interface {
	PushBottom(task Task)
	IsEmpty() bool //returns whether the queue is empty
	PopTop() Task
	PopBottom() Task
	Size() int
}

/******** DO NOT MODIFY ANY OF THE ABOVE INTERFACES/TYPES *********************/

// NewUnBoundedDEQueue returns an empty UnBoundedDEQueue
func NewUnBoundedDEQueue() DEQueue {
	return &Queue{top: nil, bottom: nil, size: 0}
}

type Queue struct {
	top    *Node
	bottom *Node
	mu sync.Mutex
	size int
}

type Node struct {
	task Task
	next *Node
	prev *Node
}

// Function to Push at the bottom of the queue
func (deq *Queue) PushBottom(task Task) {

	// Ensuring that the queue is locked before doing any operation
	deq.mu.Lock()
	defer deq.mu.Unlock()
	newNode := &Node{task: task, next: nil, prev: deq.bottom}

	// The queue is empty to start with
	if deq.bottom == nil {
		deq.bottom = newNode
		deq.top = newNode
	} else {
		deq.bottom.next = newNode
		deq.bottom = newNode
	}
	deq.size += 1
}

// Function to check if the queue is Empty
func (deq *Queue) IsEmpty() bool {
	deq.mu.Lock()
	defer deq.mu.Unlock()
	if deq.top == nil && deq.bottom == nil {
		return true
	}
	return false
}

// Function to get an element from the top of the queue
func (deq *Queue) PopTop() Task {
	deq.mu.Lock()
	defer deq.mu.Unlock()
	task := deq.top.task

	// There is only one item in the queue
	if deq.top.next == nil {
		deq.bottom = nil
		deq.top = nil
	} else {
		deq.top = deq.top.next
		deq.top.prev = nil
	}
	deq.size -= 1
	return task
}

// Function to pop from the bottom of the queue
func (deq *Queue) PopBottom() Task {
	deq.mu.Lock()
	defer deq.mu.Unlock()
	task := deq.bottom.task
	if deq.bottom.prev == nil {
		deq.bottom = nil
		deq.top = nil
	} else {
		deq.bottom = deq.bottom.prev
		deq.bottom.next = nil
	}
	deq.size -= 1
	return task
}

// Function to get the size of the queue
func (deq *Queue) Size() int {
	return deq.size
}