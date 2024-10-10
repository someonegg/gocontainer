// Copyright 2024 someonegg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heap_test

import (
	"fmt"
	"github.com/someonegg/gocontainer/cmp"
	"github.com/someonegg/gocontainer/heap"
)

func ExampleHeap() {
	h := heap.New[int](6, cmp.Less[int])

	fmt.Println(h.Len())

	h.Push(3)
	h.Push(4)
	h.Push(5)
	h.Push(1)
	h.Push(6)
	h.Push(0)
	h.Push(7)
	h.Push(2)
	h.Push(1)
	h.Push(10)

	fmt.Println(h.Len())

	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())

	fmt.Println(h.Len())

	// Output:
	// 0
	// 10
	// 0
	// 1
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 10
	// 0
}

func ExampleFixedHeap() {
	h := heap.NewFixed[int](5, cmp.Less[int])

	fmt.Println(h.Len())

	h.Push(3)
	h.Push(4)
	h.Push(5)
	h.Push(1)
	fmt.Println(h.Push(6))
	fmt.Println(h.Push(0))
	fmt.Println(h.Push(7))
	fmt.Println(h.Push(2))
	fmt.Println(h.Push(1))
	fmt.Println(h.Push(10))

	fmt.Println(h.Len())

	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())

	fmt.Println(h.Len())

	// Output:
	// 0
	// 0 false
	// 0 true
	// 1 true
	// 2 true
	// 1 true
	// 3 true
	// 5
	// 4
	// 5
	// 6
	// 7
	// 10
	// 0
}
