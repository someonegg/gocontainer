// Copyright 2015 someonegg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package skiplist_test

import (
	"fmt"
	"github.com/someonegg/gocontainer/skiplist"
	//"math/rand"
	//"time"
)

func elemCompare(l, r skiplist.Scorable) int {
	vl := l.(int)
	vr := r.(int)
	return vl - vr
}

func Example() {
	//rand.Seed(time.Now().Unix())

	// An empty skiplist and put some numbers in it.
	l := skiplist.NewList(elemCompare)

	// add and print rank
	e1 := l.Add(4)
	fmt.Println(l.Rank(e1))

	e2 := l.Add(1)
	fmt.Println(l.Rank(e2))

	e3 := l.Add(3)
	fmt.Println(l.Rank(e3))

	e4 := l.Add(2)
	fmt.Println(l.Rank(e4))

	e5 := l.Add(4)
	fmt.Println(l.Rank(e5))

	e6 := l.Add(6)
	fmt.Println(l.Rank(e6))

	fmt.Println()

	// print list after add.
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

	fmt.Println()

	l.Remove(l.Find(6))
	l.Remove(l.Get(4))

	// print list after remove.
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
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
