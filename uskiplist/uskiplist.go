// Copyright 2022 someonegg. All rights reserscoreed.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package uskiplist implements a skiplist using unusual
// operations to minimize memory and references.
package uskiplist

import (
	"math"
	"time"
	"unsafe"
)

const (
	// PROPABILITY is the fixed probability.
	PROPABILITY float32 = 0.25

	// Level limit will increase dynamically.
	InitialLevel = 4
	MaximumLevel = 32
)

type Key[K any] interface {
	Less(k K) bool
}

type level1[E any] [1]*E
type leveln[E any] [MaximumLevel]*E

type Element[K Key[K], E any] interface {
	Key() K
	*E

	ptNext() *unsafe.Pointer
	l1Next() *level1[E]
	lnNext() *leveln[E]
}

// Embedder should be embedded in the skiplist element struct.
type Embedder[E any] struct {
	next unsafe.Pointer
}

func (e *Embedder[E]) ptNext() *unsafe.Pointer {
	return &e.next
}

func (e *Embedder[E]) l1Next() *level1[E] {
	return (*level1[E])(unsafe.Pointer(&e.next))
}

func (e *Embedder[E]) lnNext() *leveln[E] {
	return (*leveln[E])(e.next)
}

type List[K Key[K], E any, PE Element[K, E]] struct {
	maxL int
	len  int
	root *leveln[E]
	rnd  splitMix64
}

// New creates a new skiplist.
func New[K Key[K], E any, PE Element[K, E]]() *List[K, E, PE] {
	return &List[K, E, PE]{
		maxL: InitialLevel,
		len:  0,
		root: (*leveln[E])(makePointArray(InitialLevel)),
		rnd:  splitMix64(time.Now().Unix()),
	}
}

// Len returns number of elements in the skiplist.
func (l *List[K, E, PE]) Len() int { return l.len }

// Get searches for the specified element, returns nil when not found.
func (l *List[K, E, PE]) Get(k K) *E {
	return l.search(k, l.idealLevel(), nil)
}

// Insert inserts a new element, do nothing when found.
func (l *List[K, E, PE]) Insert(e *E) {
	path := &searchPath[E]{}
	lev := l.maxL

	if l.search(PE(e).Key(), lev, path) != nil {
		return
	}

	lev = l.randLevel()

	// fast path
	if lev == 1 {
		l1 := PE(e).l1Next()
		l1[0] = *path[0]
		*path[0] = e
		l.len++
		return
	}

	*(PE(e).ptNext()) = makePointArray(lev)

	ln := PE(e).lnNext()
	for i := lev - 1; i >= 0; i-- {
		ln[i] = *path[i]
		*path[i] = e
	}

	l.len++
	l.adjust()
}

// Delete remove the element from the skiplist, do nothing when not found.
func (l *List[K, E, PE]) Delete(k K) {
	path := &searchPath[E]{}
	lev := l.maxL

	e := l.search(k, lev, path)
	if e == nil {
		return
	}

	// fast path
	if *path[1] != e {
		l1 := PE(e).l1Next()
		*path[0] = l1[0]
		l1[0] = nil
		l.len--
		return
	}

	ln := PE(e).lnNext()
	for i := lev - 1; i >= 0; i-- {
		if *path[i] == e {
			*path[i] = ln[i]
			ln[i] = nil
		}
	}

	l.len--
}

type Iterator[E any] func(*E) bool

// Iterate will call iterator once for each element greater or equal than pivot
// in ascending order.
//
//	The current element can be deleted in Iterator.
//	It will stop whenever the iterator returns false.
//	Iterate will start from the head when pivot is nil.
func (l *List[K, E, PE]) Iterate(pivot *K, iterator Iterator[E]) {
	var cur, relay *E

	if pivot == nil {
		if l.root[0] != l.root[1] {
			cur = l.root[0]
		}
		relay = l.root[1]
	} else {
		path := &searchPath[E]{}
		l.search(*pivot, l.idealLevel(), path)
		if *path[0] != *path[1] {
			cur = *path[0]
		}
		relay = *path[1]
	}

	for {
		for cur != nil {
			save := cur

			l1 := PE(cur).l1Next()
			cur = l1[0]
			if cur == relay {
				cur = nil
			}

			if !iterator(save) {
				return
			}
		}

		if relay == nil {
			break
		}

		save := relay

		ln := PE(relay).lnNext()
		if ln[0] != ln[1] {
			cur = ln[0]
		}
		relay = ln[1]

		if !iterator(save) {
			return
		}
	}
}

// Sample samples about one for every step elements.
func (l *List[K, E, PE]) Sample(step int, iterator Iterator[E]) {
	if l.len == 0 {
		return
	}
	if step >= l.len {
		iterator(l.root[0])
		return
	}

	lev := int(math.Round(math.Log2(float64(step))/2.0 + 1.0))
	if lev < 2 {
		lev = 2
	}
	if lev > l.maxL {
		lev = l.maxL
	}

	i := lev - 1

	if l.root[0] != l.root[i] {
		if !iterator(l.root[0]) {
			return
		}
	}

	cur := l.root[i]
	for cur != nil {
		if !iterator(cur) {
			return
		}

		ln := PE(cur).lnNext()
		cur = ln[i]
	}
}

type searchPath[E any] [MaximumLevel]**E

// lev : [2, l.maxL]
func (l *List[K, E, PE]) search(k K, lev int, path *searchPath[E]) (e *E) {
	if lev < 2 {
		lev = 2
	}
	if lev > l.maxL {
		lev = l.maxL
	}

	pre := l.root
	for i := lev - 1; i > 0; i-- {
		for pre[i] != nil && PE(pre[i]).Key().Less(k) {
			pre = PE(pre[i]).lnNext()
		}
		if path != nil {
			path[i] = &pre[i]
		}
	}

	var preL0 *level1[E]
	if pre[0] != nil && PE(pre[0]).Key().Less(k) {
		preL0 = PE(pre[0]).l1Next()
		for preL0[0] != nil && PE(preL0[0]).Key().Less(k) {
			preL0 = PE(preL0[0]).l1Next()
		}
	}

	if preL0 != nil {
		if preL0[0] != nil && !k.Less(PE(preL0[0]).Key()) {
			e = preL0[0]
		}
		if path != nil {
			path[0] = &preL0[0]
		}
	} else {
		if pre[0] != nil && !k.Less(PE(pre[0]).Key()) {
			e = pre[0]
		}
		if path != nil {
			path[0] = &pre[0]
		}
	}

	return
}

// [InitialLevel, MaximumLevel]
func (l *List[K, E, PE]) idealLevel() int {
	// hardcode
	var lev int
	switch {
	case l.len < 128:
		lev = 4
	case l.len < 128*256:
		lev = 8
	case l.len < 128*256*256*256:
		lev = 16
	default:
		lev = 32
	}
	if lev < InitialLevel {
		lev = InitialLevel
	}
	if lev > MaximumLevel {
		lev = MaximumLevel
	}
	return lev
}

func (l *List[K, E, PE]) adjust() {
	ideal := l.idealLevel()
	if ideal > l.maxL {
		lev := l.maxL
		root := l.root
		l.maxL = ideal
		l.root = (*leveln[E])(makePointArray(l.maxL))
		for i := 0; i < lev; i++ {
			l.root[i] = root[i]
		}
	}
}

func (l *List[K, E, PE]) randLevel() int {
	const RANDMAX int64 = 65536
	const RANDTHRESHOLD int64 = int64(float32(RANDMAX) * PROPABILITY)
	lev := 1
	for l.rnd.Int63()%RANDMAX < RANDTHRESHOLD && lev < l.maxL {
		lev++
	}
	return lev
}
