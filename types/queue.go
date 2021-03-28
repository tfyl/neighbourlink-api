package types

import "sync"


type qnode struct {
	data interface{}   // data payload
	next *qnode        // next node in queue
}


type Queue struct {
	head *qnode         // head of queue
	tail *qnode         // tail of queue
	l int               // length of queue
	mutex *sync.Mutex   // mutex to control queue
}


//initialise queue and return queue class
func NewQueue() *Queue {
	// instantiates the queue struct / object
	q := &Queue{}
	// assigns the attribute a mutex
	q.mutex = &sync.Mutex{}
	return q
}

// Method of Queue
// push new item to queue (add to the tail/end)
func (q *Queue) Push (i interface{}){
	// lock the mutex to block other processes from running and interfering
	q.mutex.Lock()
	// defer the lock until the function is returned
	defer q.mutex.Unlock()

	// instantiate a new node with a nil pointer (null) for the next node
	n := &qnode{
		data: i,
		next: nil,
	}

	// if tail is nil the node must be the first item
	// thus the node is both the tail and head
	if q.tail == nil{
		q.tail = n
		q.head = n
	} else {
		// assign the next value for the current tail
		q.tail.next = n
		// change the tail to the new node
		q.tail = n
	}

	// increase length of queue
	q.l ++
}

// Method of Queue
// pop item off queue (from the front)
func (q *Queue) Pop () (i interface{}){
	// lock the mutex to block other processes from running and interfering
	q.mutex.Lock()
	// defer the lock until the function is returned
	defer q.mutex.Unlock()

	if q.head == nil {
		return nil
	}

	i = q.head.data
	q.head = q.head.next

	q.l --
	// reduce queue length

	return
}

// Method of Queue
// peek item on queue (check first value without changing the data structure)
func (q *Queue) Peek () (i interface{}){
	// lock the mutex to block other processes from running and interfering
	q.mutex.Lock()
	// defer the lock until the function is returned
	defer q.mutex.Unlock()

	n := q.head
	if n == nil {
		return nil
	}

	return n.data
}

// Method of Queue
// return length of data structure
func (q *Queue) Len () int {
	// lock the mutex to block other processes from running and interfering
	q.mutex.Lock()
	// defer the lock until the function is returned
	defer q.mutex.Unlock()
	return q.l
}
