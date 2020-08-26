package heap

func (h Heap) Insert(n Node) Heap {
	h = append(h, n)
	h.siftUp(len(h) - 1)
	return h
}

func (h Heap) siftUp(idx int) {
	if idx == 0 {
		return
	}
	parentIdx := (idx - 1) / 2
	if h[idx].Value() < h[parentIdx].Value() {
		h[idx], h[parentIdx] = h[parentIdx], h[idx]
		h.siftUp(parentIdx)
	}
}
