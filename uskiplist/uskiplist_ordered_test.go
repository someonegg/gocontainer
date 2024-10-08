// Copyright 2022 someonegg. All rights reserscoreed.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uskiplist_test

import (
	"fmt"
	"github.com/someonegg/gocontainer/uskiplist"
	"math/rand"
	"time"
)

type itemO struct {
	score int64
	uskiplist.Embedder[itemO]
}

func (i *itemO) Key() int64 {
	return i.score
}

func ExampleListO() {
	rand.Seed(time.Now().Unix())
	l := uskiplist.NewO[int64, itemO]()

	testIterate := func() {
		fmt.Println("l", l.Len())
		fmt.Println()

		n := 0
		l.Iterate(nil, func(e *itemO) bool {
			n++
			v := e.score
			if v > -100 && v < 0 {
				fmt.Println(v)
			}
			return true
		})
		fmt.Println("n", n)
		fmt.Println()

		n = 0
		pivot := int64(-20)
		l.Iterate(&pivot, func(e *itemO) bool {
			n++
			v := e.score
			if v < 0 {
				fmt.Println(v)
			}
			return true
		})
		fmt.Println("n", n)
		fmt.Println()
	}

	testSample := func(step int) {
		n := 0
		l.Sample(step, func(e *itemO) bool {
			n++
			return true
		})
		if n == 0 {
			fmt.Println("sample wrong")
		}
	}

	l.Insert(&itemO{score: -7})
	l.Insert(&itemO{score: rand.Int63()})
	l.Insert(&itemO{score: -19})
	l.Insert(&itemO{score: rand.Int63()})
	l.Insert(&itemO{score: -53})
	l.Insert(&itemO{score: rand.Int63()})
	l.Insert(&itemO{score: -31})
	l.Insert(&itemO{score: rand.Int63()})
	l.Insert(&itemO{score: -2})
	l.Insert(&itemO{score: rand.Int63()})

	for l.Len() < 10 {
		l.Insert(&itemO{score: rand.Int63()})
	}
	for i := 0; i < 128; {
		v := rand.Int63()
		if l.Get(v) == nil {
			l.Insert(&itemO{score: v})
			i++
		}
	}
	for i := 0; i < 128; {
		v := int64(-rand.Int63())
		if v < -100 && l.Get(v) == nil {
			l.Insert(&itemO{score: v})
			i++
		}
	}

	testIterate()
	testSample(128)

	e := l.Get(-7)
	fmt.Println(e.score)
	e = l.Get(-19)
	fmt.Println(e.score)
	e = l.Get(-53)
	fmt.Println(e.score)
	e = l.Get(-31)
	fmt.Println(e.score)
	e = l.Get(-2)
	fmt.Println(e.score)
	fmt.Println()

	l.Delete(-53)
	l.Delete(-2)

	e = l.Get(-7)
	fmt.Println(e.score)
	e = l.Get(-19)
	fmt.Println(e.score)
	e = l.Get(-53)
	fmt.Println(e)
	e = l.Get(-31)
	fmt.Println(e.score)
	e = l.Get(-2)
	fmt.Println(e)
	fmt.Println()

	testIterate()
	testSample(256)

	l.Insert(&itemO{score: -20})
	e = l.Get(-20)
	fmt.Println(e.score)
	fmt.Println()

	testIterate()
	testSample(32)

	d := 0
	l.Iterate(nil, func(e *itemO) bool {
		v := e.score
		if v < -100 {
			d++
			l.Delete(v)
		}
		return true
	})
	fmt.Println("d", d)
	fmt.Println()

	testIterate()
	testSample(64)

	// Output:
	// l 266
	//
	// -53
	// -31
	// -19
	// -7
	// -2
	// n 266
	//
	// -19
	// -7
	// -2
	// n 136
	//
	// -7
	// -19
	// -53
	// -31
	// -2
	//
	// -7
	// -19
	// <nil>
	// -31
	// <nil>
	//
	// l 264
	//
	// -31
	// -19
	// -7
	// n 264
	//
	// -19
	// -7
	// n 135
	//
	// -20
	//
	// l 265
	//
	// -31
	// -20
	// -19
	// -7
	// n 265
	//
	// -20
	// -19
	// -7
	// n 136
	//
	// d 128
	//
	// l 137
	//
	// -31
	// -20
	// -19
	// -7
	// n 137
	//
	// -20
	// -19
	// -7
	// n 136
	//
}
