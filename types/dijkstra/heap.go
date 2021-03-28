package dijkstra

import HeapPackage "container/heap"

// class/struct for path , contains a value and node attribute
type path struct {
	value int
	nodes []string
}
// class/struct that stores an array of paths (object)
type minPath []path

// methods of the minPath class
//
// returns length of the object array
func (h minPath) Len() int           { return len(h) }
// compares two values and returns bool if it is less
func (h minPath) Less(i, j int) bool { return h[i].value < h[j].value }
// swaps the position of two items in the array
func (h minPath) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
// push a new item to the array
func (h *minPath) Push(x interface{}) {
	*h = append(*h, x.(path))
}
// pop item off array
func (h *minPath) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}



// class / struct that defines a heap , it stores a pointer to a minPath struct
type heap struct {
	values *minPath
}
// creates new heap
func newHeap() *heap {
	return &heap{values: &minPath{}}
}
// pushes item to heap
func (h *heap) push(p path) {
	HeapPackage.Push(h.values, p)
}
// pop off heap
func (h *heap) pop() path {
	i := HeapPackage.Pop(h.values)
	return i.(path)
}