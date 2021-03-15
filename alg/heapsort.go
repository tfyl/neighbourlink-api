package alg

import "fmt"

type heap struct {
	array []*Hnode
}

func NewHeap(arr []*Hnode) *heap {
	return &heap{
		array: arr,
	}

}

type Hnode struct {
	Value float64
	Data interface{}
}

func (h *heap) lChildIndex (i int) int {  // return index of left child
	return 2*i
}

func (h *heap) lChildValue (i int) float64 {  // return value of left child
	return h.array[2*i].Value
}

func (h *heap) rChildIndex (i int) int {  // return index of right child
	return 2*i + 1
}

func (h *heap) rChildValue (i int) float64 {  // return value of right child
	return h.array[2*i+1].Value
}

func (h *heap) swap(a, b int){  // swap two nodes/items around
	temp := h.array[a]
	h.array[a] = h.array[b]
	h.array[b] = temp
}

func (h *heap) leaf(i int, size int) bool { // check of current index is greater than height and i is less than or equal to the size
	if i >= (size/2) && i <= size {
		return true
	}
	return false
}

func (h *heap) heapifyDown(current int, size int) { // heapify down function that is carried out recursively until the array is is a max heap
	if h.leaf(current, size) {
		return
	}
	smallest := current

	leftChildIndex := h.lChildIndex(current)
	rightRightIndex := h.rChildIndex(current)

	if leftChildIndex < size && h.lChildValue(current) < h.array[smallest].Value {
		smallest = leftChildIndex
	}
	if rightRightIndex < size && h.rChildValue(current) < h.array[smallest].Value {
		smallest = rightRightIndex
	}

	if smallest != current {
		h.swap(current, smallest)
		h.heapifyDown(smallest, size)
	}

	return
}

func (h *heap) buildHeap(size int) {  // create heap
	for index := (size / 2) - 1; index >= 0; index-- {
		h.heapifyDown(index, size)
	}
}

func (h *heap) Sort () {  // sorts tree and moves head to end once sorted
	size := len(h.array)
	h.buildHeap(size)
	for i := size - 1; i > 0; i-- {
		// Move current head to end of array
		h.swap(0, i)
		h.heapifyDown(0, i)
	}
}

func (h *heap) Print() { // prints array
	for _, val := range h.array {
		fmt.Println(val)
	}
}

func (h *heap) ReturnArray() []*Hnode {  // returns sorted array
	return h.array
}
