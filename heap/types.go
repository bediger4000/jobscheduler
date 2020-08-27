package heap

// Node is the interface tupe of items
// stored in the heap
type Node interface {
	Value() int64
	IsNil() bool
	String() string
}

// Heap as an array: simplifies maintaining the shape
// of the heap.
type Heap []Node

/*
Thus the children of the node at position n would
2n + 1 and 2n + 2 in a zero-based array.
Computing the index of the parent node of n-th element is also
straightforward.
Similarly, for zero-based arrays, is the parent is
located at position (n-1)/2 (floored).
*/
