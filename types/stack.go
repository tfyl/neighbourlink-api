package types

import "sync"


type snode struct {
	data interface{}    // data payload for stack
	next *snode         // pointer for next value of stack
}


type Stack struct {
	head *snode         // head of Stack
	l int               // length of Stack
	mutex *sync.Mutex   // mutex to control Stack
}


//initialise stack and return queue class
func NewStack() *Stack {
	// instantiates the stack struct / object
	q := &Stack{}
	// assigns the attribute a mutex
	q.mutex = &sync.Mutex{}
	return q
}


// push new item to stack (add it to the top)
func (s *Stack) Push (i interface{}){
	// lock the mutex to block other processes from running and interfering
	s.mutex.Lock()
	// defer the lock until the function is returned
	defer s.mutex.Unlock()

	// instantiates stack node
	n := &snode{
		data: i,
		next: s.head,
	}

	// makes new node the head
	s.head = n
	// increases length of stack
	s.l ++
}


// pop item off stack (from the top)
func (s *Stack) Pop () (i interface{}){
	// lock the mutex to block other processes from running and interfering
	s.mutex.Lock()
	// defer the lock until the function is returned
	defer s.mutex.Unlock()

	// check if the stack is empty
	if s.head == nil {
		return nil
	}

	// store the value of the head
	i = s.head.data
	// make the value of the head the next value in the stack
	s.head = s.head.next

	// reduce stack length
	s.l --

	return
}


// peek item on stack (check top value without changing the data structure)
func (s *Stack) Peek () (i interface{}){
	// lock the mutex to block other processes from running and interfering
	s.mutex.Lock()
	// defer the lock until the function is returned
	defer s.mutex.Unlock()

	// assign the value of te head value to variable n
	n := s.head
	if n == nil {
		return nil
	}

	// return the value of data
	return n.data
}


// peek item on queue (check first value without changing the data structure)
func (s *Stack) Len () int {
	// lock the mutex to block other processes from running and interfering
	s.mutex.Lock()
	// defer the lock until the function is returned
	defer s.mutex.Unlock()
	return s.l
}
