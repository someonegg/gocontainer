// Copyright 2022 someonegg. All rights reserscoreed.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package uskiplist implements a skiplist using unsafe
// operations to minimize memory and references.
package uskiplist

import (
	"math"
	"math/rand"
	"sync"
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

type lockedSource struct {
	lk  sync.Mutex
	src rand.Source
}

func (r *lockedSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *lockedSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}

var gRandSource rand.Source = &lockedSource{src: rand.NewSource(time.Now().Unix())}

// Header should be embedded in the skiplist element struct:
//
// generically
//
//	type Element struct {
//	    uskiplist.Header
//	    ...
//	}
//
// sort by string
//
//	type Element struct {
//	    uskiplist.Header
//	    abc string
//	    ...
//	}
//
// sort by integer
//
//	type Element struct {
//	    uskiplist.Header
//	    abc int64
//	    ...
//	}
//
// sort by float
//
//	type Element struct {
//	    uskiplist.Header
//	    abc float64
//	    ...
//	}
type Header struct {
	next unsafe.Pointer
}

type levels [MaximumLevel]*element

type element struct {
	next *levels
}

type elementL0 struct {
	next *elementL0
}

type LessFunc func(l, r unsafe.Pointer) bool

type List struct {
	less LessFunc
	maxL int
	len  int
	root element
}

// New creates a new skiplist.
func New(less LessFunc) *List {
	if less == nil {
		panic("less is nil")
	}

	l := &List{
		less: less,
		maxL: InitialLevel,
		len:  0,
	}
	l.root.next = (*levels)(makePointArray(l.maxL))

	return l
}

// Len returns number of elements in the skiplist.
func (l *List) Len() int { return l.len }

// Get searches for the specified element, returns nil when not found.
func (l *List) Get(e unsafe.Pointer) unsafe.Pointer {
	t, _, _ := l.search(e, l.idealLevel(), nil)
	return t
}

// Insert inserts a new element, do nothing when found.
func (l *List) Insert(e unsafe.Pointer) {
	path := &searchPath{}
	lev := l.maxL

	t, _, _ := l.search(e, lev, path)
	if t != nil {
		return
	}

	t = e
	lev = l.randLevel()

	// fast path
	if lev == 1 {
		tt := (*elementL0)(t)
		if path.preL0 != nil {
			tt.next = path.preL0.next
			path.preL0.next = tt
		} else {
			tt.next = (*elementL0)(unsafe.Pointer(path.pre[0].next[0]))
			path.pre[0].next[0] = (*element)(unsafe.Pointer(tt))
		}
		l.len++
		return
	}

	tt := (*element)(t)
	tt.next = (*levels)(makePointArray(lev))

	for i := lev - 1; i > 0; i-- {
		tt.next[i] = path.pre[i].next[i]
		path.pre[i].next[i] = tt
	}

	if path.preL0 != nil {
		tt.next[0] = (*element)(unsafe.Pointer(path.preL0.next))
		path.preL0.next = (*elementL0)(unsafe.Pointer(tt))
	} else {
		tt.next[0] = path.pre[0].next[0]
		path.pre[0].next[0] = tt
	}

	l.len++
	l.adjust()
}

// Delete remove the element from the skiplist, do nothing when not found.
func (l *List) Delete(e unsafe.Pointer) {
	path := &searchPath{}
	lev := l.maxL

	t, _, _ := l.search(e, lev, path)
	if t == nil {
		return
	}

	// fast path
	if unsafe.Pointer(path.pre[1].next[1]) != t {
		tt := (*elementL0)(t)
		next := tt.next
		tt.next = nil
		if path.preL0 != nil {
			path.preL0.next = next
		} else {
			path.pre[0].next[0] = (*element)(unsafe.Pointer(next))
		}
		l.len--
		return
	}

	tt := (*element)(t)

	for i := lev - 1; i > 0; i-- {
		if path.pre[i].next[i] == tt {
			next := tt.next[i]
			tt.next[i] = nil
			path.pre[i].next[i] = next
		}
	}

	next := tt.next[0]
	tt.next[0] = nil
	if path.preL0 != nil {
		path.preL0.next = (*elementL0)(unsafe.Pointer(next))
	} else {
		path.pre[0].next[0] = next
	}

	l.len--
}

type Iterator func(e unsafe.Pointer) bool

// Iterate will call iterator once for each element greater or equal than pivot
// in ascending order.
//
//	The current element can be deleted in Iterator.
//	It will stop whenever the iterator returns false.
//	Iterate will start from the head when pivot is nil.
func (l *List) Iterate(pivot unsafe.Pointer, iterator Iterator) {
	var (
		cur   *element
		curL0 *elementL0
	)

	if pivot == nil {
		cur = &l.root
		if cur.next[0] != cur.next[1] {
			curL0 = (*elementL0)(unsafe.Pointer(cur.next[0]))
		}
	} else {
		_, cur, curL0 = l.search(pivot, l.idealLevel(), nil)
		if curL0 != nil {
			curL0 = curL0.next
			if unsafe.Pointer(curL0) == unsafe.Pointer(cur.next[1]) {
				curL0 = nil
			}
		} else {
			if cur.next[0] != cur.next[1] {
				curL0 = (*elementL0)(unsafe.Pointer(cur.next[0]))
			}
		}
	}

	next := cur.next[1]

	for {
		for curL0 != nil {
			save := curL0

			curL0 = curL0.next
			if unsafe.Pointer(curL0) == unsafe.Pointer(next) {
				curL0 = nil
			}

			if !iterator(unsafe.Pointer(save)) {
				return
			}
		}

		if next == nil {
			break
		}

		cur = next
		next = cur.next[1]

		if cur.next[0] != next {
			curL0 = (*elementL0)(unsafe.Pointer(cur.next[0]))
		}

		if !iterator(unsafe.Pointer(cur)) {
			return
		}

	}
}

// Sample samples about one for every step elements.
func (l *List) Sample(step int, iterator Iterator) {
	if l.len == 0 {
		return
	}
	if step >= l.len {
		iterator(unsafe.Pointer(l.root.next[0]))
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

	if l.root.next[0] != l.root.next[i] {
		if !iterator(unsafe.Pointer(l.root.next[0])) {
			return
		}
	}

	cur := l.root.next[i]
	for cur != nil {
		if !iterator(unsafe.Pointer(cur)) {
			return
		}
		cur = cur.next[i]
	}
}

// [InitialLevel, MaximumLevel]
func (l *List) idealLevel() int {
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

func (l *List) adjust() {
	ideal := l.idealLevel()
	if ideal > l.maxL {
		lev := l.maxL
		next := l.root.next
		l.maxL = ideal
		l.root.next = (*levels)(makePointArray(l.maxL))
		for i := 0; i < lev; i++ {
			l.root.next[i] = next[i]
		}
	}
}

func (l *List) randLevel() int {
	const RANDMAX int64 = 65536
	const RANDTHRESHOLD int64 = int64(float32(RANDMAX) * PROPABILITY)
	lev := 1
	for gRandSource.Int63()%RANDMAX < RANDTHRESHOLD && lev < l.maxL {
		lev++
	}
	return lev
}

type searchPath struct {
	pre   levels
	preL0 *elementL0
}

// lev : [2, l.maxL]
func (l *List) search(e unsafe.Pointer, lev int, path *searchPath) (t unsafe.Pointer, pre *element, preL0 *elementL0) {
	if lev < 2 {
		lev = 2
	}
	if lev > l.maxL {
		lev = l.maxL
	}

	pre = &l.root
	for i := lev - 1; i > 0; i-- {
		for pre.next[i] != nil && l.less(unsafe.Pointer(pre.next[i]), e) {
			pre = pre.next[i]
		}
		if path != nil {
			path.pre[i] = pre
		}
	}

	preL0 = nil
	if pre.next[0] != nil && l.less(unsafe.Pointer(pre.next[0]), e) {
		preL0 = (*elementL0)(unsafe.Pointer(pre.next[0]))
		for preL0.next != nil && l.less(unsafe.Pointer(preL0.next), e) {
			preL0 = preL0.next
		}
	}

	if preL0 != nil {
		if preL0.next != nil && !l.less(e, unsafe.Pointer(preL0.next)) {
			t = unsafe.Pointer(preL0.next)
		}
	} else {
		if pre.next[0] != nil && !l.less(e, unsafe.Pointer(pre.next[0])) {
			t = unsafe.Pointer(pre.next[0])
		}
	}

	if path != nil {
		path.pre[0] = pre
		path.preL0 = preL0
	}

	return
}
