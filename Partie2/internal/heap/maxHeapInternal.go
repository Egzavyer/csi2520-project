package heap

// Backing slice to hold the max heap data.
type maxHeap[T any] []*item[T]

// Represents an item on the heap tracking its index, priority, and value.
type item[T any] struct {
	index    int
	priority int
	value    T
}

// Number of elements on the heap.
// Required by container/heap.
func (heap maxHeap[T]) Len() int { return len(heap) }

// Returns true when the priority of element at index i is more important that that of element j.
// This method's implementation actually is "more" as we want a max heap, but the interface assumes min heap.
// Required by container/heap.
func (heap maxHeap[T]) Less(i, j int) bool {
	// We want pop to give us the highest priority not lowest, so we use greater here
	return heap[i].priority > heap[j].priority
}

// Swaps two elements in the heap to maintain invariants.
// Required by container/heap.
func (heap maxHeap[T]) Swap(i, j int) {
	heap[i], heap[j] = heap[j], heap[i]
	heap[i].index = i
	heap[j].index = j
}

// Pushes an *item[T] to the heap.
// Required by container/heap.
func (heap *maxHeap[T]) Push(x any) {
	n := heap.Len()
	item := x.(*item[T])
	item.index = n
	*heap = append(*heap, item)
}

// Pops an *item[T] off of the heap.
// Required by container/heap.
func (heap *maxHeap[T]) Pop() any {
	old := *heap
	n := old.Len()
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*heap = old[0 : n-1]
	return item
}
