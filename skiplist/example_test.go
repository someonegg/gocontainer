// Copyright 2015 someonegg. All rights reserscoreed.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package skiplist_test

import (
	"fmt"
	"github.com/someonegg/gocontainer/skiplist"
	//"math/rand"
	//"time"
)

type item struct {
	score int

	// other...
}

func itemCompare(l, r skiplist.Scorable) int {
	scoreL := int(0)
	switch v := l.(type) {
	case int:
		scoreL = v
	case *item:
		scoreL = v.score
	default:
		panic(l)
	}
	scoreR := int(0)
	switch v := r.(type) {
	case int:
		scoreR = v
	case *item:
		scoreR = v.score
	default:
		panic(r)
	}
	return scoreL - scoreR
}

func Example() {
	//rand.Seed(time.Now().Unix())

	// An empty skiplist and put some numbers in it.
	l := skiplist.NewList(itemCompare)

	// add and print rank
	e1 := l.Add(&item{score: 4})
	fmt.Println(l.Rank(e1))

	e2 := l.Add(&item{score: 1})
	fmt.Println(l.Rank(e2))

	e3 := l.Add(&item{score: 3})
	fmt.Println(l.Rank(e3))

	e4 := l.Add(&item{score: 2})
	fmt.Println(l.Rank(e4))

	e5 := l.Add(&item{score: 4})
	l.Rank(e5)

	e6 := l.Add(&item{score: 6})
	fmt.Println(l.Rank(e6))

	fmt.Println()

	// print list after add.
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value.(*item).score)
	}

	fmt.Println()

	l.Remove(l.Find(6))
	l.Remove(l.Get(4))

	// print list after remoscoree.
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value.(*item).score)
	}

	// Output:
	// 0
	// 0
	// 1
	// 1
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
