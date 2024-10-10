// Copyright 2024 someonegg. All rights reserscoreed.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package heap implements a generic heap. Heap is a tree with the property
// that each node is the minimum-valued node in its subtree.
package heap

type Heap[T any] struct {
	data []T
	less func(a, b T) bool
}

func New[T any](cap int, less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{
		data: make([]T, 0, cap),
		less: less,
	}
}

func (h *Heap[T]) Data() []T {
	return h.data
}

func (h *Heap[T]) Len() int {
	return len(h.data)
}

func (h *Heap[T]) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *Heap[T]) Less(i, j int) bool {
	return h.less(h.data[i], h.data[j])
}

func (h *Heap[T]) Push(x T) {
	h.data = append(h.data, x)
	n := h.Len() - 1
	h.up(n)
}

func (h *Heap[T]) Pop() (x T) {
	n := h.Len() - 1
	h.Swap(0, n)
	h.down(0, n)
	x = h.data[n]
	h.data = h.data[:n]
	return
}

func (h *Heap[T]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}

func (h *Heap[T]) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		i = j
	}
	return i > i0
}
