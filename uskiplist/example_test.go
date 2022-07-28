// Copyright 2022 someonegg. All rights reserscoreed.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uskiplist_test

import (
	"fmt"
	"github.com/someonegg/gocontainer/uskiplist"
	"math/rand"
	"time"
	"unsafe"
)

type item struct {
	uskiplist.Header
	score int64
}

func Example() {
	rand.Seed(time.Now().Unix())
	l := uskiplist.NewByInteger()

	testIterate := func() {
		fmt.Println("l", l.Len())
		fmt.Println()

		n := 0
		l.Iterate(nil, func(e unsafe.Pointer) bool {
			n++
			v := (*item)(e).score
			if v > -100 && v < 0 {
				fmt.Println(v)
			}
			return true
		})
		fmt.Println("n", n)
		fmt.Println()

		n = 0
		l.IterateByInteger(-20, func(e unsafe.Pointer) bool {
			n++
			v := (*item)(e).score
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
		l.Sample(step, func(e unsafe.Pointer) bool {
			n++
			return true
		})
		if n == 0 {
			fmt.Println("sample wrong")
		}
	}

	l.Insert(unsafe.Pointer(&item{score: -7}))
	l.Insert(unsafe.Pointer(&item{score: rand.Int63()}))
	l.Insert(unsafe.Pointer(&item{score: -19}))
	l.Insert(unsafe.Pointer(&item{score: rand.Int63()}))
	l.Insert(unsafe.Pointer(&item{score: -53}))
	l.Insert(unsafe.Pointer(&item{score: rand.Int63()}))
	l.Insert(unsafe.Pointer(&item{score: -31}))
	l.Insert(unsafe.Pointer(&item{score: rand.Int63()}))
	l.Insert(unsafe.Pointer(&item{score: -2}))
	l.Insert(unsafe.Pointer(&item{score: rand.Int63()}))

	for l.Len() < 10 {
		l.Insert(unsafe.Pointer(&item{score: rand.Int63()}))
	}
	for i := 0; i < 128; {
		v := rand.Int63()
		if l.GetByInteger(v) == nil {
			l.Insert(unsafe.Pointer(&item{score: v}))
			i++
		}
	}
	for i := 0; i < 128; {
		v := -rand.Int63()
		if v < -100 && l.GetByInteger(v) == nil {
			l.Insert(unsafe.Pointer(&item{score: v}))
			i++
		}
	}

	testIterate()
	testSample(128)

	e := l.GetByInteger(-7)
	fmt.Println((*item)(e).score)
	e = l.GetByInteger(-19)
	fmt.Println((*item)(e).score)
	e = l.GetByInteger(-53)
	fmt.Println((*item)(e).score)
	e = l.GetByInteger(-31)
	fmt.Println((*item)(e).score)
	e = l.GetByInteger(-2)
	fmt.Println((*item)(e).score)
	fmt.Println()

	l.DeleteByInteger(-53)
	l.DeleteByInteger(-2)

	e = l.GetByInteger(-7)
	fmt.Println((*item)(e).score)
	e = l.GetByInteger(-19)
	fmt.Println((*item)(e).score)
	e = l.GetByInteger(-53)
	fmt.Println(e)
	e = l.GetByInteger(-31)
	fmt.Println((*item)(e).score)
	e = l.GetByInteger(-2)
	fmt.Println(e)
	fmt.Println()

	testIterate()
	testSample(256)

	l.Insert(unsafe.Pointer(&item{score: -20}))
	e = l.GetByInteger(-20)
	fmt.Println((*item)(e).score)
	fmt.Println()

	testIterate()
	testSample(32)

	d := 0
	l.Iterate(nil, func(e unsafe.Pointer) bool {
		v := (*item)(e).score
		if v < -100 {
			d++
			l.Delete(e)
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
