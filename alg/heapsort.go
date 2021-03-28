package alg

import "fmt"

type heap struct {
	array []*Hnode
}

func NewHeap(arr []*Hnode) *heap {
	// Creates new heap object
	return &heap{
		array: arr,
	}

}

type Hnode struct {
	// The "Class" definition for each Heap Node
	Value float64
	Data interface{}
}

// method of heap
// return index of left child
func (h *heap) lChildIndex (i int) int {
	return 2*i
}

// method of heap
// return value of left child
func (h *heap) lChildValue (i int) float64 {
	return h.array[2*i].Value
}

// method of heap
// return index of right child
func (h *heap) rChildIndex (i int) int {
	return 2*i + 1
}

// method of heap
// return value of right child
func (h *heap) rChildValue (i int) float64 {
	return h.array[2*i+1].Value
}

// method of heap
// swap two nodes/items around
func (h *heap) swap(a, b int){
	temp := h.array[a]
	h.array[a] = h.array[b]
	h.array[b] = temp
}

// method of heap
// check of current index is greater than height and i is less than or equal to the size
func (h *heap) leaf(i int, size int) bool {
	if i >= (size/2) && i <= size {
		return true
	}
	return false
}

// method of heap
// heapify down function that is carried out recursively until the array is is a max heap
func (h *heap) heapifyDown(current int, size int) {
	if h.leaf(current, size) {
		return
	}
	smallest := current

	// gets index of the left child
	leftChildIndex := h.lChildIndex(current)
	// gets index of the right child
	rightRightIndex := h.rChildIndex(current)

	// compares the values of the left child to the smallest in array
	if leftChildIndex < size && h.lChildValue(current) < h.array[smallest].Value {
		smallest = leftChildIndex
	}

	// compares the values of the right child to the smallest in array
	if rightRightIndex < size && h.rChildValue(current) < h.array[smallest].Value {
		smallest = rightRightIndex
	}

	// if it hasn't reached the smallest node it hasn't finished
	if smallest != current {
		// swaps the current and the smallest node
		h.swap(current, smallest)
		// calls the same function (recursion)
		h.heapifyDown(smallest, size)
	}

	return
}

// method of heap
// create heap
func (h *heap) buildHeap(size int) {
	for index := (size / 2) - 1; index >= 0; index-- {
		h.heapifyDown(index, size)
	}
}

// method of heap
// sorts tree and moves head to end once sorted
func (h *heap) Sort () {
	size := len(h.array)
	h.buildHeap(size)
	for i := size - 1; i > 0; i-- {
		// Move current head to end of array
		h.swap(0, i)
		h.heapifyDown(0, i)
	}
}

// method of heap
// prints array
func (h *heap) Print() {
	for _, val := range h.array {
		fmt.Println(val)
	}
}

// method of heap
// returns array
func (h *heap) ReturnArray() []*Hnode {
	return h.array
}
