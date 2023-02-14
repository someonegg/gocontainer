package uskiplist

import (
	"math"
	"unsafe"
)

const (
	// PROPABILITY is the fixed probability.
	PROPABILITY float32 = 0.25

	// Level limit will increase dynamically.
	InitialLevel = 4
	MaximumLevel = 32
)

type Comparable[T any] interface {
	Less(v T) bool
}

type Element[K Comparable[K], E any] interface {
	GetKey() K
	header() *ElementHeader
	*E
}

type ElementHeader struct {
	next unsafe.Pointer
}

func (h *ElementHeader) header() *ElementHeader {
	return h
}

type levels [MaximumLevel]*element

type element struct {
	next *levels
}

type elementL0 struct {
	next *elementL0
}

type searchPath struct {
	pre [MaximumLevel]*element
}

type List[K Comparable[K], E any, PE Element[K, E]] struct {
	maxL int
	len  int
	root E
}

func NewList[K Comparable[K], E any, PE Element[K, E]]() *List[K, E, PE] {
	l := &List[K, E, PE]{
		maxL: InitialLevel,
		len:  0,
	}
	(PE(&l.root)).header().next = makePointArray(InitialLevel)
	return l
}

func (l *List[K, E, PE]) Len() int { return l.len }

func (l *List[K, E, PE]) Get(k K) *E {
	e, _, _ := l.search(k, l.idealLevel(), nil)
	return e
}

func (l *List[K, E, PE]) Insert(newE *E) {
	lev := l.maxL
	path := searchPath{}

	e, preL0, _ := l.search((PE(newE)).GetKey(), lev, &path)

	if e != nil {
		return
	}

	lev = l.randLevel()

	var pNextL0 unsafe.Pointer
	eHeader := (*element)(unsafe.Pointer(newE))

	if lev != 1 {
		eHeader.next = (*levels)(makePointArray(lev))
		for i := lev - 1; i > 0; i-- {
			eHeader.next[i] = path.pre[i].next[i]
			path.pre[i].next[i] = (*element)(unsafe.Pointer(newE))
		}
		pNextL0 = unsafe.Pointer(&(eHeader.next[0]))
	} else {
		pNextL0 = unsafe.Pointer(&(eHeader.next))
	}

	// process level 0
	if preL0 != nil {
		*(*unsafe.Pointer)(pNextL0) = unsafe.Pointer(preL0.next)
		preL0.next = (*elementL0)(unsafe.Pointer(eHeader))
	} else {
		*(*unsafe.Pointer)(pNextL0) = unsafe.Pointer(path.pre[0].next[0])
		path.pre[0].next[0] = eHeader
	}

	l.len++
	l.adjust()
}

func (l *List[K, E, PE]) Delete(k K) {
	lev := l.maxL
	path := searchPath{}

	e, preL0, _ := l.search(k, lev, &path)
	if e == nil {
		return
	}

	eHeader := (*element)(unsafe.Pointer(e))
	var pNextL0 unsafe.Pointer

	if path.pre[1].next[1] == eHeader {
		for i := lev - 1; i > 0; i-- {
			if path.pre[i].next[i] == eHeader {
				next := eHeader.next[i]
				eHeader.next[i] = nil
				path.pre[i].next[i] = next
			}
		}
		pNextL0 = unsafe.Pointer(&(eHeader.next[0]))
	} else {
		pNextL0 = unsafe.Pointer(&(eHeader.next))
	}

	next := *(*unsafe.Pointer)(pNextL0)
	*(*unsafe.Pointer)(pNextL0) = nil
	if preL0 != nil {
		preL0.next = (*elementL0)(next)
	} else {
		path.pre[0].next[0] = (*element)(next)
	}

	l.len--

}

type Iterator[K Comparable[K], E any, PE Element[K, E]] func(*E) bool

// Iterate will call iterator once for each element greater or equal than pivot
// in ascending order.
//
//	The current element can be deleted in Iterator.
//	It will stop whenever the iterator returns false.
//	Iterate will start from the head when pivot is nil.
func (l *List[K, E, PE]) Iterate(iterator Iterator[K, E, PE]) {
	var (
		curL0 *elementL0
		cur   = (*element)(unsafe.Pointer(&l.root))
	)
	if cur.next[0] != cur.next[1] {
		curL0 = (*elementL0)(unsafe.Pointer(cur.next[0]))
	}

	l.iterate(cur, curL0, iterator)
}

func (l *List[K, E, PE]) IterateWithPivot(pivot K, iterator Iterator[K, E, PE]) {
	_, curL0, cur := l.search(pivot, l.idealLevel(), nil)
	if curL0 == nil {
		if cur.next[0] != cur.next[1] {
			curL0 = (*elementL0)(unsafe.Pointer(cur.next[0]))
		}
	} else {
		if unsafe.Pointer(curL0.next) != unsafe.Pointer(cur.next[1]) {
			curL0 = curL0.next
		} else {
			curL0 = nil
		}
	}
	l.iterate(cur, curL0, iterator)
}

func (l *List[K, E, PE]) iterate(cur *element, curL0 *elementL0, iterator Iterator[K, E, PE]) {
	next := cur.next[1]
	for {
		for curL0 != nil {
			save := curL0
			curL0 = curL0.next

			if unsafe.Pointer(curL0) == unsafe.Pointer(next) {
				curL0 = nil
			}

			if !iterator((*E)(unsafe.Pointer(save))) {
				return
			}
		}

		if next == nil {
			return
		}

		cur = next
		next = cur.next[1]

		if cur.next[0] != next {
			curL0 = (*elementL0)(unsafe.Pointer(cur.next[0]))
		}

		if !iterator((*E)(unsafe.Pointer(cur))) {
			return
		}
	}

}

func (l *List[K, E, PE]) Sample(step int, iterator Iterator[K, E, PE]) {
	if l.len == 0 {
		return
	}
	if step >= l.len {
		iterator(&l.root)
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

	rootHeader := (*element)(unsafe.Pointer(&l.root))

	if rootHeader.next[0] != rootHeader.next[i] {
		if !iterator((*E)(unsafe.Pointer(rootHeader.next[0]))) {
			return
		}
	}

	cur := rootHeader.next[i]
	for cur != nil {
		if !iterator((*E)(unsafe.Pointer(cur))) {
			return
		}
		cur = cur.next[i]
	}
}

// search return the last node that smaller than target
func (l *List[K, E, PE]) search(key K, lev int, path *searchPath) (e *E, preL0 *elementL0, pre *element) {
	if lev < 2 {
		lev = 2
	}

	if lev > l.maxL {
		lev = l.maxL
	}

	pre = (*element)(unsafe.Pointer((PE(&l.root)).header()))

	for i := lev - 1; i > 0; i-- {
		for pre.next[i] != nil && l.less(unsafe.Pointer(pre.next[i]), key) {
			pre = pre.next[i]
		}
		if path != nil {
			path.pre[i] = pre
		}
	}

	preL0 = nil

	if pre.next[0] != nil && l.less(unsafe.Pointer(pre.next[0]), key) {
		preL0 = (*elementL0)(unsafe.Pointer(pre.next[0]))
		for preL0.next != nil && l.less(unsafe.Pointer(preL0.next), key) {
			preL0 = preL0.next
		}
	}

	if preL0 != nil {
		if preL0.next != nil && !key.Less((PE((*E)(unsafe.Pointer(preL0.next)))).GetKey()) {
			e = (*E)(unsafe.Pointer(preL0.next))
		}
	} else {
		if pre.next[0] != nil && !key.Less((PE((*E)(unsafe.Pointer(pre.next[0])))).GetKey()) {
			e = (*E)(unsafe.Pointer(pre.next[0]))
		}
	}

	if path != nil {
		path.pre[0] = pre
	}

	return
}

func (l *List[K, E, PE]) less(elem unsafe.Pointer, key K) bool {
	return (PE((*E)(elem))).GetKey().Less(key)
}

func (l *List[K, E, PE]) randLevel() int {
	const RANDMAX int64 = 65536
	const RANDTHRESHOLD int64 = int64(float32(RANDMAX) * PROPABILITY)
	lev := 1
	for gRandSource.Int63()%RANDMAX < RANDTHRESHOLD && lev < l.maxL {
		lev++
	}
	return lev
}

// [InitialLevel, MaximumLevel]
func (l *List[K, E, PE]) idealLevel() int {
	//
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
		root := (*element)(unsafe.Pointer((PE(&l.root)).header()))
		next := root.next
		l.maxL = ideal
		root.next = (*levels)(makePointArray(l.maxL))
		for i := 0; i < lev; i++ {
			root.next[i] = next[i]
		}
	}
}
