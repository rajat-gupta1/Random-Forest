package concurrent


import (
	"math/rand"
	"time"
	"sync"
)

// Struct for implementing work balance executor
type WorkBalancingExecutor struct {
	capacity  int
	threshold int
	thresholdBalance int
	deq       []DEQueue
	threadVal int
	globalQ DEQueue
	done int
	mu sync.Mutex
}

// NewWorkBalancingExecutor returns an ExecutorService that is implemented using the work-balancing algorithm.
// @param capacity - The number of goroutines in the pool
// @param threshold - The number of items that a goroutine in the pool can
// grab from the executor in one time period. For example, if threshold = 10
// this means that a goroutine can grab 10 items from the executor all at
// once to place into their local queue before grabbing more items. It's
// not required that you use this parameter in your implementation.
// @param thresholdBalance - The threshold used to know when to perform
// balancing. Remember, if two local queues are to be balanced the
// difference in the sizes of the queues must be greater than or equal to
// thresholdBalance. You must use this parameter in your implementation.
func NewWorkBalancingExecutor(capacity, threshold, thresholdBalance int) ExecutorService {
	var deq []DEQueue

	// Setting the number of local queues equal to the number of threads
	for i := 0; i < capacity; i++ {
		deq = append(deq, NewUnBoundedDEQueue())
	}
	globalQ := NewUnBoundedDEQueue()
	thisExecutor := &WorkBalancingExecutor{capacity: capacity, threshold: threshold, deq: deq, threadVal: 0, globalQ: globalQ, done: 0, thresholdBalance: thresholdBalance}
	
	// Setting threads to run the DistributeTask function
	for i := 0; i < capacity; i++ {
		go DistributeTask(thisExecutor, i)
	}
	return thisExecutor
}

// Function to push work in golabal queue
func (WSE *WorkBalancingExecutor) Submit(task interface{}) Future {
	c := make(chan interface{}, 1)
	thisTask := taskSt{task: task, fut: NewFuture(c)}
	WSE.globalQ.PushBottom(thisTask)

	// Returning future
	return thisTask.fut
}

// Function to do the task
func DoTask (thisTask interface{}, ThreadId int) {
	thisTask2, _ := thisTask.(taskSt)
	task3, ok := thisTask2.task.(Runnable)
	if ok {
		task3.Run()
		thisTask2.fut.c <- nil
	} else {
		task2, _ := thisTask2.task.(Callable)
		result := task2.Call()
		thisTask2.fut.c <- result
	}
}

// Function to Get the task 
func GetTask(WSE *WorkBalancingExecutor, ThreadId int, done int) int {
	ctr := 0
	for {
		if WSE.globalQ.IsEmpty() {

			// Either we are done or we have some work to do
			if done == 1 || WSE.deq[ThreadId].Size() > 0{
				break
			} else {

				// Take one task from any other thread to start the proceedings
				for i := 0; i < WSE.capacity - 1; i++ {
					if WSE.deq[(i + ThreadId + 1) % WSE.capacity].Size() > 0 {
						ctr = 1
						break
					}
				} 

				// Task taken from any thread
				if ctr == 1 {
					ctr = 0
					break
				}
			}

			// Sleep for some time to avoid unnecessary spinning
			time.Sleep(time.Duration(rand.Int31n(300)) * time.Millisecond)

		} else {
			// There are some tasks in the global queue. Locking to ensure that 
			// no other thread take the task in the meantime
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

// Function to help in WorkDistribution
func DistributeWork(WSE *WorkBalancingExecutor, ThreadId int) {
	// Generate random number and see if balancing is to be done
	size := WSE.deq[ThreadId].Size()
	randNum := rand.Intn(size + 1)
	if randNum == size {
		victim := rand.Intn(WSE.capacity)
		size2 := WSE.deq[victim].Size()

		// If size of victim is higher than size of the current thread
		if size2 > size {
			if size2 - size > WSE.thresholdBalance {
				for WSE.deq[victim].Size() - WSE.deq[ThreadId].Size() > 0 {
					WSE.deq[ThreadId].PushBottom(WSE.deq[victim].PopTop())
				}
			}
		} else {
			if size - size2 > WSE.thresholdBalance {
				for WSE.deq[ThreadId].Size() - WSE.deq[victim].Size() > 0 {
					WSE.deq[victim].PushBottom(WSE.deq[ThreadId].PopTop())
				}
			}
		}
	}
}

// function to Distribute Tasks
func DistributeTask(WSE *WorkBalancingExecutor, ThreadId int) {
	done := 0
	for {
		done = GetTask(WSE, ThreadId, done)
		
		// Breaking the loop in case Shutdown has been called and all tasks have been completed
		if done == 1 {
			break
		}

		// Do tasks till they are done
		for WSE.deq[ThreadId].Size() > 0 {
			
			thisTask := WSE.deq[ThreadId].PopBottom()
			DoTask(thisTask, ThreadId)
			DistributeWork(WSE, ThreadId)
			
		}
	}
}

// Function to signal if there is no more work
func (WSE *WorkBalancingExecutor) Shutdown() {
	WSE.done = 1
}