package types

import (
	"math"
	"sync"
)


type pqnode struct {
	priority int
	data interface{}
	next *pqnode
}


type PQueue struct {
	head *pqnode          // head of queue
	l int               // length of queue
	mutex *sync.Mutex   // mutex to control queue
}


func NewPQueue() *PQueue {   //initialise priority queue and return queue class
	q := &PQueue{head:&pqnode{priority:math.MinInt64,next:&pqnode{priority:math.MinInt64}}}
	q.mutex = &sync.Mutex{}
	return q
}


func (q *PQueue) Push (i interface{},p int){   // push new item to queue (add to the tail/end)
	q.l ++
	q.mutex.Lock()
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
		if c.next != nil && c.next.priority > n.priority {  // checks if the priority of the item is less than the next item, if so insert
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


func (q *PQueue) Pop () (i interface{}){  // pop item off priority queue (from the front)
	q.mutex.Lock()
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


func (q *PQueue) Peek () (i interface{}){  // peek item from priority queue (check first value without changing the data structure)
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.l == 0 {
		return nil
	}

	n := q.head

	return n.data
}


func (q *PQueue) Len () int {  // return length of data structure
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.l
}
