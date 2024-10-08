// Copyright 2022 someonegg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uskiplist_test

import (
	"fmt"
	"github.com/someonegg/gocontainer/uskiplist"
	"math/rand"
	"time"
)

type item struct {
	uskiplist.Embedder[item]

	k string
	v int64
}

func (i *item) Key() string {
	return i.k
}

func Example() {
	rand.Seed(time.Now().Unix())
	l := uskiplist.NewO[string, item]()

	l.Insert(&item{k: "a", v: 1})
	l.Insert(&item{k: "b", v: 2})
	l.Insert(&item{k: "c", v: 3})
	fmt.Println(l.Len())

	e := l.Get("b")
	fmt.Println(e.v)

	l.Delete("b")
	fmt.Println(l.Len())

	e = l.Get("b")
	fmt.Println(e)

	e = l.Get("d")
	fmt.Println(e)

	l.Insert(&item{k: "d", v: 4})
	fmt.Println(l.Len())

	e = l.Get("d")
	fmt.Println(e.v)

	l.Iterate(nil, func(e *item) bool {
		fmt.Println(e.v)
		return true
	})

	pivot := "c"
	l.Iterate(&pivot, func(e *item) bool {
		fmt.Println(e.v)
		return true
	})

	// Output:
	// 3
	// 2
	// 2
	// <nil>
	// <nil>
	// 3
	// 4
	// 1
	// 3
	// 4
	// 3
	// 4
}
