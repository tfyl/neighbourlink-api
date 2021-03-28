package types

import (
	"math"
	"sync"
)


type pqnode struct {
	priority int        // priority from which it is sorted
	data interface{}    // data payload that is sorted
	next *pqnode        // stores pointer to the next priority queue node
}


type PQueue struct {
	head *pqnode        // head of queue
	l int               // length of queue
	mutex *sync.Mutex   // mutex to control queue
}


//initialise priority queue and return queue class
func NewPQueue() *PQueue {
	q := &PQueue{head:&pqnode{priority:math.MinInt64,next:&pqnode{priority:math.MinInt64}}}
	q.mutex = &sync.Mutex{}
	return q
}

// Method of PQueue
// push new item to queue (add to the tail/end)
func (q *PQueue) Push (i interface{},p int){
	q.l ++
	// lock the mutex to block other processes from running and interfering
	q.mutex.Lock()
	// defer the lock until the function is returned
	defer q.mutex.Unlock()

	n := &pqnode{
		data: i,
		next: nil,
		priority:p,
	}


	// acts as a cursor searching through the queue
	c := q.head


	// length has already increased by one so must start from i=-1
	for i := -1; i < q.l; i++ {
		 // checks if the priority of the item is less than the next item, if so insert
		if c.next != nil && c.next.priority > n.priority {
			n.next = c.next
			c.next = n
			break

		} else {
			if c.next == nil {
				c.next = n
			}
			c = c.next
			continue
		}
	}
	return
}

// Method of PQueue
// pop item off priority queue (from the front)
func (q *PQueue) Pop () (i interface{}){
	// lock the mutex to block other processes from running and interfering
	q.mutex.Lock()
	// defer the lock until the function is returned
	defer q.mutex.Unlock()

	if q.l == 0 {
		return nil
	}
	//
	i = q.head.next.next.data
	q.head = q.head.next

	q.l --
	// reduce queue length

	return
}

// Method of PQueue
// peek item from priority queue (check first value without changing the data structure)
func (q *PQueue) Peek () (i interface{}){
	// lock the mutex to block other processes from running and interfering
	q.mutex.Lock()
	// defer the lock until the function is returned
	defer q.mutex.Unlock()

	if q.l == 0 {
		return nil
	}

	n := q.head

	return n.data
}

// Method of PQueue
// return length of data structure
func (q *PQueue) Len () int {
	// lock the mutex to block other processes from running and interfering
	q.mutex.Lock()
	// defer the lock until the function is returned
	defer q.mutex.Unlock()

	return q.l
}
