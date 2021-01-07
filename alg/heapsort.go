package alg

import "fmt"

type heap struct {
	array []*hnode
}

func NewHeap(arr []*hnode) *heap {
	return &heap{
		array: arr,
	}

}

type hnode struct {
	value int
	data interface{}
}

func (h *heap) lChildIndex (i int) int {
	return 2*i
}

func (h *heap) lChildValue (i int) int {
	return h.array[2*i].value
}

func (h *heap) rChildIndex (i int) int {
	return 2*i + 1
}

func (h *heap) rChildValue (i int) int {
	return h.array[2*i+1].value
}

func (h *heap) swap(a, b int){
	temp := h.array[a]
	h.array[a] = h.array[b]
	h.array[b] = temp
}

func (h *heap) leaf(i int, size int) bool {
	if i >= (size/2) && i <= size {
		return true
	}
	return false
}

func (h *heap) heapifyUp(current int, size int) {
	if h.leaf(current, size) {
		return
	}
	smallest := current

	leftChildIndex := h.lChildIndex(current)
	rightRightIndex := h.rChildIndex(current)

	if leftChildIndex < size && h.lChildValue(current) < h.array[smallest].value {
		smallest = leftChildIndex
	}
	if rightRightIndex < size && h.rChildValue(current) < h.array[smallest].value {
		smallest = rightRightIndex
	}

	if smallest != current {
		h.swap(current, smallest)
		h.heapifyUp(smallest, size)
	}

	return
}

func (h *heap) buildHeap(size int) {
	for index := (size / 2) - 1; index >= 0; index-- {
		h.heapifyUp(index, size)
	}
}

func (h *heap) sort () {
	size := len(h.array)
	h.buildHeap(size)
	for i := size - 1; i > 0; i-- {
		// Move current head to end of array
		h.swap(0, i)
		h.heapifyUp(0, i)
	}
}

func (h *heap) print() {
	for _, val := range h.array {
		fmt.Println(val)
	}
}

func (h *heap) returnArray() []*hnode {
	return h.array
}
