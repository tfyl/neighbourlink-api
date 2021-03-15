package types

import "sync"


type snode struct {
	data interface{}
	next *snode
}


type Stack struct {
	head *snode          // head of Stack
	l int               // length of Stack
	mutex *sync.Mutex   // mutex to control Stack
}


func NewStack() *Stack {     //initialise stack and return queue class
	q := &Stack{}
	q.mutex = &sync.Mutex{}
	return q
}


func (s *Stack) Push (i interface{}){    // push new item to stack (add it to the top)
	s.mutex.Lock()
	defer s.mutex.Unlock()

	n := &snode{
		data: i,
		next: s.head,
	}

	s.head = n
	s.l ++
}


func (s *Stack) Pop () (i interface{}){ // pop item off stack (from the top)
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.head == nil {
		return nil
	}

	i = s.head.data
	s.head = s.head.next

	s.l --
	// reduce stack length

	return
}


func (s *Stack) Peek () (i interface{}){  // peek item on stack (check top value without changing the data structure)
	s.mutex.Lock()
	defer s.mutex.Unlock()

	n := s.head
	if n == nil {
		return nil
	}

	return n.data
}


func (s *Stack) Len () int {  // peek item on queue (check first value without changing the data structure)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.l
}
