// Copyright 2015 someonegg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package skiplist implements a skip list. Compared with the classical
// version, there are two changes:
//
//	this implementation allows for repeated elements.
//	there is a back pointer, so it's a doubly linked list.
//
// List will be sorted by score:
//
//	in ascending order.
//	with rank(0-based), also in ascending order.
//	"what is the score, how to compare" is defined by the user.
package skiplist

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	// PROPABILITY is the fixed probability.
	PROPABILITY float32 = 0.25

	DefaultLevel = 16
	MaximumLevel = 32
)

// Scorable object can be passed to CompareFunc.
type Scorable interface{}

// CompareFunc can compare two scorable objects, returns
//
//	<0 if l <  r
//	 0 if l == r
//	>0 if l >  r
type CompareFunc func(l, r Scorable) int

// Element is an element of a skip list.
type Element struct {
	// The value stored with this element.
	Value Scorable

	lev  []level
	list *List
}

type level struct {
	next *Element
	prev *Element
	span int
}

// Next returns the next list element or nil.
func (e *Element) Next() *Element {
	if p := e.next(); e.list != nil && p != e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *Element) Prev() *Element {
	if p := e.prev(); e.list != nil && p != e.list.root {
		return p
	}
	return nil
}

func (e *Element) next() *Element {
	return e.lev[0].next
}

func (e *Element) prev() *Element {
	return e.lev[0].prev
}

// List represents a skip list.
type List struct {
	maxL int
	rndS rand.Source
	comp CompareFunc
	root *Element
	len  int
}

// NewList creates a new skip list, with DefaultLevel\compare.
func NewList(compare CompareFunc) *List {
	return NewListEx(DefaultLevel, compare)
}

// NewListEx creates a new skip list, with maxLevel\compare.
func NewListEx(maxLevel int, compare CompareFunc) *List {
	if maxLevel < 1 || maxLevel > MaximumLevel {
		panic("maxLevel < 1 or maxLevel > MaximumLevel")
	}
	if compare == nil {
		panic("compare is nil")
	}

	l := &List{
		maxL: maxLevel,
		rndS: rand.NewSource(time.Now().Unix()),
		comp: compare,
		root: &Element{
			lev:  make([]level, maxLevel),
			list: nil,
		},
	}

	for i := 0; i < l.maxL; i++ {
		l.root.lev[i].next = l.root
		l.root.lev[i].prev = l.root
		l.root.lev[i].span = 0
	}

	return l
}

// Len returns the number of elements of list l. The complexity is O(1).
func (l *List) Len() int { return l.len }

// Front returns the first element of list l or nil.
func (l *List) Front() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.next()
}

// Back returns the last element of list l or nil.
func (l *List) Back() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.prev()
}

// Get the element at rank, return nil if rank is invalid.
//
//	0 <= valid rank < list.Len()
func (l *List) Get(rank int) *Element {
	if rank < 0 || rank >= l.len {
		return nil
	}

	e, found := l.searchToRank(rank, nil)
	if !found || e == l.root {
		panic("impossible")
	}

	return e
}

// Find the first element equal to score, return nil if not found.
// If there are multiple elements equal to score, you can use the
// "Element" to traverse them.
func (l *List) Find(score Scorable) *Element {
	if score == nil {
		return nil
	}

	e, found := l.searchToScore(score, nil)
	if found && e == l.root {
		panic("impossible")
	}

	if !found {
		return nil
	}
	return e
}

// Rank will calculate current rank of the element, return -1 if not in the list.
func (l *List) Rank(e *Element) int {
	if e.list != l {
		return -1
	}

	path := &searchPath{}
	l.searchPathOf(e, path)

	span := 0
	for _, v := range path.levSpan {
		span += v
	}

	return span - 1
}

// Add an element to the list.
func (l *List) Add(v Scorable) *Element {
	e := &Element{Value: v}
	l.add(e)
	return e
}

func (l *List) add(e *Element) {
	path := &searchPath{}

	ee, found := l.searchToScore(e.Value, path)
	if found && ee == l.root {
		panic("impossible")
	}

	randON := true

	// repeated element
	if found {
		if len(ee.lev) == 1 {
			// only 1 level, insert before.
			path.prev[0] = path.prev[0].prev()
			path.levSpan[0]--
		} else {
			// more than 1 level, force 1 level on newer.
			randON = false
		}
	}

	nlev := 1
	if randON {
		nlev = l.randLevel()
	}

	//fmt.Println(nlev, randON)

	e.lev = make([]level, nlev)

	revspan := 0
	for i := 0; i < nlev; i++ {
		p := path.prev[i]
		n := p.lev[i].next
		p.lev[i].next = e
		e.lev[i].prev = p
		e.lev[i].next = n
		n.lev[i].prev = e

		e.lev[i].span = p.lev[i].span - revspan
		p.lev[i].span = revspan + 1
		revspan += path.levSpan[i]
	}

	for i := nlev; i < l.maxL; i++ {
		path.prev[i].lev[i].span++
	}

	e.list = l
	l.len++
}

func (l *List) randLevel() int {
	const RANDMAX int64 = 65536
	const RANDTHRESHOLD int64 = int64(float32(RANDMAX) * PROPABILITY)
	nlev := 1
	for l.rndS.Int63()%RANDMAX < RANDTHRESHOLD && nlev < l.maxL {
		nlev++
	}
	return nlev
}

// Remove an element from the list.
func (l *List) Remove(e *Element) {
	if e.list != l {
		return
	}
	l.remove(e)
}

func (l *List) remove(e *Element) {
	path := &searchPath{}
	l.searchPathOf(e, path)

	for i := 0; i < len(e.lev); i++ {
		n := e.lev[i].next
		p := e.lev[i].prev
		p.lev[i].next = n
		n.lev[i].prev = p
		p.lev[i].span += e.lev[i].span - 1
	}

	for i := len(e.lev); i < l.maxL; i++ {
		path.prev[i].lev[i].span--
	}

	e.lev = nil
	e.list = nil
	l.len--
}

// searchPath represents search path of skip list.
type searchPath struct {
	prev    [MaximumLevel]*Element
	levSpan [MaximumLevel]int
}

func (l *List) searchPathOf(e *Element, path *searchPath) {
	path.prev[0] = e
	path.levSpan[0] = 0

	ilev := 0
	for {
		idle, levSpan := true, 0

		for i := ilev + 1; i < len(e.lev); i++ {
			idle = false
			path.prev[i] = e
			path.levSpan[i] = 0
			ilev++
		}

		for e != l.root && (ilev+1) >= len(e.lev) {
			idle = false
			e = e.lev[ilev].prev
			levSpan += e.lev[ilev].span
		}

		if idle {
			break
		}

		path.levSpan[ilev] = levSpan
	}
}

// return
//
//	<0 goto down
//	 0 found
//	>0 goto next
type poscompFunc func(ilev int, p, n *Element) int

func (l *List) searchToPos(poscomp poscompFunc, path *searchPath) (*Element, bool) {
	found := false

	p := l.root

	for i := l.maxL - 1; i >= 0; i-- {
		levSpan := 0

		equal := false
		n := p.lev[i].next
		for n != l.root {
			ret := poscomp(i, p, n)
			if ret < 0 {
				break
			}
			levSpan += p.lev[i].span
			p = n
			n = n.lev[i].next
			if ret == 0 {
				equal = true
				break
			}
		}

		if path != nil {
			path.prev[i] = p
			path.levSpan[i] = levSpan
		}

		if equal {
			found = true
			for i--; path != nil && i >= 0; i-- {
				path.prev[i] = p
				path.levSpan[i] = 0
			}
			break
		}
	}

	return p, found
}

// searchToXXX will find the element that is closest to XXX.
// If the "path" is not nil, it will be filled.

func (l *List) searchToScore(score Scorable, path *searchPath) (*Element, bool) {
	poscomp := func(ilev int, p, n *Element) int {
		return l.comp(score, n.Value)
	}

	return l.searchToPos(poscomp, path)
}

func (l *List) searchToRank(rank int, path *searchPath) (*Element, bool) {
	span := rank + 1
	poscomp := func(ilev int, p, n *Element) int {
		ret := span - p.lev[ilev].span
		if ret >= 0 {
			span = ret
		}
		return ret
	}
	return l.searchToPos(poscomp, path)
}

func (l *List) dump() {
	fmt.Println("TotalLevel:", l.maxL, " ", "Length:", l.len)
	fmt.Println()
	for i := l.maxL - 1; i >= 0; i-- {
		fmt.Println("Level:", i)
		e := l.root

		fmt.Print("  ")
		for {
			fmt.Println("\t", e.lev[i].span)

			e = e.lev[i].next
			if e == l.root {
				break
			}
			fmt.Print(e, " ")
		}
		fmt.Println()
	}
}
