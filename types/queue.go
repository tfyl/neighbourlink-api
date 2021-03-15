package types

import "sync"


type qnode struct {
	data interface{}
	next *qnode
}


type Queue struct {
	head *qnode          // head of queue
	tail *qnode          // tail of queue
	l int               // length of queue
	mutex *sync.Mutex   // mutex to control queue
}


func NewQueue() *Queue {    //initialise queue and return queue class
	q := &Queue{}
	q.mutex = &sync.Mutex{}
	return q
}


func (q *Queue) Push (i interface{}){  // push new item to queue (add to the tail/end)
	q.mutex.Lock()
	defer q.mutex.Unlock()

	n := &qnode{
		data: i,
		next: nil,
	}

	if q.tail == nil{
		q.tail = n
		q.head = n
	} else {
		q.tail.next = n
		q.tail = n
	}

	q.l ++
}


func (q *Queue) Pop () (i interface{}){  // pop item off queue (from the front)
	q.mutex.Lock()
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


func (q *Queue) Peek () (i interface{}){  // peek item on queue (check first value without changing the data structure)
	q.mutex.Lock()
	defer q.mutex.Unlock()

	n := q.head
	if n == nil {
		return nil
	}

	return n.data
}


func (q *Queue) Len () int { // return length of data structure
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.l
}
