// Copyright 2015 someonegg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package skiplist_test

import (
	"fmt"
	"github.com/someonegg/gocontainer/skiplist"
)

type elem struct {
	skiplist.Node
	v int
}

func elemCompare(l, r skiplist.Scorable) int {
	el := l.(*elem)
	er := r.(*elem)
	return el.v - er.v
}

func Example() {
	// An empty skiplist and put some numbers in it.
	l := skiplist.NewList(elemCompare)

	// add and print rank
	fmt.Println(l.Add(&elem{v: 4}))
	fmt.Println(l.Add(&elem{v: 1}))
	fmt.Println(l.Add(&elem{v: 3}))
	fmt.Println(l.Add(&elem{v: 2}))
	fmt.Println(l.Add(&elem{v: 4}))
	e6 := &elem{v: 6}
	fmt.Println(l.Add(e6))

	fmt.Println()

	// print list after add.
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.(*elem).v)
	}

	fmt.Println()

	l.Remove(e6)
	l.Remove(l.FindOfRank(4))

	// print list after remove.
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.(*elem).v)
	}

	// Output:
	// 0
	// 0
	// 1
	// 1
	// 3
	// 5
	//
	// 1
	// 2
	// 3
	// 4
	// 4
	// 6
	//
	// 1
	// 2
	// 3
	// 4
}
