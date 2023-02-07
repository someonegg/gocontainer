// Copyright 2022 someonegg. All rights reserscoreed.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uskiplist

import (
	"fmt"
	"math/rand"
	"time"
)

type Item struct {
	Value int64
}

func (i *Item) GetKey() ListKey {
	return ListKey(i.Value)
}

type ListKey int64

func (i ListKey) Less(i2 ListKey) bool {
	return i < i2
}

func Example() {
	rand.Seed(time.Now().Unix())
	l := NewList[ListKey, *Item]()

	testIterate := func() {
		fmt.Println("l", l.Len())
		fmt.Println()

		n := 0
		l.Iterate(func(i *Item) bool {
			n++
			if i.Value > -100 && i.Value < 0 {
				fmt.Println(i.Value)
			}
			return true
		})
		fmt.Println("n", n)
		fmt.Println()

		n = 0
		l.IterateWithPivot(-20, func(i *Item) bool {
			n++
			if i.Value < 0 {
				fmt.Println(i.Value)
			}
			return true
		})
		fmt.Println("n", n)
		fmt.Println()
	}

	testSample := func(step int) {
		n := 0
		l.Sample(step, func(i *Item) bool {
			n++
			return true
		})
		if n == 0 {
			fmt.Println("sample wrong")
		}
	}

	l.Insert(&Item{-7})
	l.Insert(&Item{rand.Int63()})
	l.Insert(&Item{-19})
	l.Insert(&Item{rand.Int63()})
	l.Insert(&Item{-53})
	l.Insert(&Item{rand.Int63()})
	l.Insert(&Item{-31})
	l.Insert(&Item{rand.Int63()})
	l.Insert(&Item{-2})
	l.Insert(&Item{rand.Int63()})

	var k int64
	for l.Len() < 10 {
		k = rand.Int63()
		l.Insert(&Item{Value: k})
	}
	for i := 0; i < 128; {
		k = rand.Int63()
		if l.Get(ListKey(k)) == nil {
			l.Insert(&Item{Value: k})
			i++
		}
	}
	for i := 0; i < 128; {
		k = -rand.Int63()
		if k < -100 {
			if l.Get(ListKey(k)) == nil {
				l.Insert(&Item{Value: k})
				i++
			}
		}
	}

	testIterate()
	testSample(128)

	e := l.Get(-7)
	fmt.Println(e.Value)
	e = l.Get(-19)
	fmt.Println(e.Value)
	e = l.Get(-53)
	fmt.Println(e.Value)
	e = l.Get(-31)
	fmt.Println(e.Value)
	e = l.Get(-2)
	fmt.Println(e.Value)
	fmt.Println()

	l.Delete(-53)
	l.Delete(-2)

	e = l.Get(-7)
	fmt.Println(e.Value)
	e = l.Get(-19)
	fmt.Println(e.Value)
	e = l.Get(-53)
	fmt.Println(e)
	e = l.Get(-31)
	fmt.Println(e.Value)
	e = l.Get(-2)
	fmt.Println(e)
	fmt.Println()

	testIterate()
	testSample(256)

	l.Insert(&Item{Value: -20})
	e = l.Get(-20)
	fmt.Println(e.Value)
	fmt.Println()

	testIterate()
	testSample(32)

	d := 0
	l.Iterate(func(i *Item) bool {
		if i.GetKey() < -100 {
			d++
			l.Delete(i.GetKey())
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
