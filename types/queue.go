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


func NewQueue() *Queue {
	q := &Queue{}
	q.mutex = &sync.Mutex{}
	return q
}


func (q *Queue) Push (i interface{}){
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


func (q *Queue) Pop () (i interface{}){
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.head == nil {
		return nil
	}

	i = q.head
	q.head = q.head.next

	q.l --
	// reduce queue length

	return
}


func (q *Queue) Peek () (i interface{}){
	q.mutex.Lock()
	defer q.mutex.Unlock()

	n := q.head
	if n == nil {
		return nil
	}

	return n.data
}


func (q *Queue) Len () int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.l
}
