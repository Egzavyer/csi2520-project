package heap

import "container/heap"

// A max heap data structure that tracks values and their corresponding priorities.
// The element with the highest priority is the only element accessible with this wrapper.
//
// Example:
//
//	heap := NewMaxHeapWithCap[string](3)
//	heap.Push(2, "cat")
//	heap.Push(5, "sheep")
//	heap.Push(9, "dog")
//
//	// What is the highest-priority element?
//	priority, animal := heap.Peek() // 9, "dog"
//
//	// Remove the highest-priority element.
//	priority, animal = heap.Pop()  // 9, "dog"
//	priority, animal = heap.Peek() // 5, "sheep"
//
// This structure re-balances itself when any modifications are made.
type MaxHeap[T any] struct {
	heap maxHeap[T]
}

// Creates a new MaxHeap[T] with a given capacity.
// Relying on default(MaxHeap[T]) instead of calling this method must be avoided.
func NewMaxHeapWithCap[T any](cap int) *MaxHeap[T] {
	h := make(maxHeap[T], 0, cap)
	return &MaxHeap[T]{h}
}

// Adds a new value to the heap with a specified priority.
func (h *MaxHeap[T]) Push(priority int, value T) {
	idx := h.heap.Len()
	item := &item[T]{
		index:    idx,
		priority: priority,
		value:    value,
	}
	heap.Push(&h.heap, item)
}

// Removes the maximum priority element from the heap and returns its priority, value.
func (h *MaxHeap[T]) Pop() (int, T) {
	elem := heap.Pop(&h.heap).(*item[T])
	return elem.priority, elem.value
}

// Checks the maximum priority element without actually modifying the heap.
// Returns -1, default(T) if the heap has no elements.
// Otherwises returns priority, value.
func (h *MaxHeap[T]) Peek() (int, T) {
	if h.Len() == 0 {
		var zero T
		return -1, zero
	}
	top := h.heap[0]
	return top.priority, top.value
}

// Gets the number of elements in the heap.
func (h *MaxHeap[T]) Len() int {
	return h.heap.Len()
}
