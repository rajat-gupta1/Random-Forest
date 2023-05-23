package concurrent

import (
	"time"
	"sync"
	"math/rand"
)

// Struct for implementing work stealing executor
type WorkStealingExecutor struct {
	capacity  int
	threshold int
	deq       []DEQueue
	threadVal int
	globalQ DEQueue
	done int
	mu sync.Mutex
}

// Creating future struct to help with get method
type Fut struct {
	c chan interface{}
}

func (fut Fut) Get() interface{} {
	return <-fut.c
}

func NewFuture(c chan interface{}) Fut {
	return Fut{c: c}
}

type taskSt struct {
	task Task
	fut Fut
}

// NewWorkStealingExecutor returns an ExecutorService that is implemented using the work-stealing algorithm.
// @param capacity - The number of goroutines in the pool
// @param threshold - The number of items that a goroutine in the pool can
// grab from the executor in one time period. For example, if threshold = 10
// this means that a goroutine can grab 10 items from the executor all at
// once to place into their local queue before grabbing more items. It's
// not required that you use this parameter in your implementation.
func NewWorkStealingExecutor(capacity, threshold int) ExecutorService {
	var deq []DEQueue

	// Setting the number of local queues equal to the number of threads
	for i := 0; i < capacity; i++ {
		deq = append(deq, NewUnBoundedDEQueue())
	}
	globalQ := NewUnBoundedDEQueue()
	thisExecutor := &WorkStealingExecutor{capacity: capacity, threshold: threshold, deq: deq, threadVal: 0, globalQ: globalQ, done: 0}
	
	// Setting threads to run the StealTask function
	for i := 0; i < capacity; i++ {
		go StealTask(thisExecutor, i)
	}
	return thisExecutor
}

// Function to push work in golabal queue
func (WSE *WorkStealingExecutor) Submit(task interface{}) Future {
	c := make(chan interface{}, 1)
	thisTask := taskSt{task: task, fut: NewFuture(c)}
	WSE.globalQ.PushBottom(thisTask)
	return thisTask.fut
}

// Function to get task
func TaskGet (WSE *WorkStealingExecutor, ThreadId int, done int) int {
	ctr := 0
	for {
		if WSE.globalQ.IsEmpty() {

			// If size of current queue = 0, steal
			if WSE.deq[ThreadId].Size() == 0 {
				for i := 0; i < WSE.capacity; i++ {
					WSE.mu.Lock()

					k := (i + ThreadId + 1) % WSE.capacity
					if WSE.deq[k].Size() > 0 {
						WSE.deq[ThreadId].PushBottom(WSE.deq[k].PopTop())
						ctr += 1
						WSE.mu.Unlock()
						break
					}
					WSE.mu.Unlock()
				} 

				// Stealing successful
				if ctr == 1 {
					ctr = 0
					break
				} else if WSE.done == 1 {
					done = 1
					break
				}

				// Sleep for some time to avoid unnecessary spinning
				time.Sleep(time.Duration(rand.Int31n(300)) * time.Millisecond)
			} else {

				// We have some tasks
				break
			}
		} else {
			WSE.mu.Lock()

			if WSE.globalQ.IsEmpty() == false {
				// Ensuring that tasks are taken only till the threshold is reached
				if WSE.deq[ThreadId].Size() < WSE.threshold {
					WSE.deq[ThreadId].PushBottom(WSE.globalQ.PopBottom())
				} else {
					WSE.mu.Unlock()
					break
				}
			}
			WSE.mu.Unlock()
		}
	}
	return done
}

// Function to Steal Task
func StealTask(WSE *WorkStealingExecutor, ThreadId int) {
	done := 0
	for {
		done = TaskGet(WSE, ThreadId, done)
		if done == 1 {
			break
		}

		for WSE.deq[ThreadId].Size() > 0 {
			thisTask := WSE.deq[ThreadId].PopBottom()
			DoTask(thisTask, ThreadId)
		}
	}
}

// Function to signal if there is no more work
func (WSE *WorkStealingExecutor) Shutdown() {
	WSE.done = 1
}