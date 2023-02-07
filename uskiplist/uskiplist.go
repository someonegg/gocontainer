package uskiplist

import (
	"math"
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

type HasKey[K any] interface {
	GetKey() K
}

type levels[K Comparable[K], PV HasKey[K]] []*element[K, PV]

type element[K Comparable[K], PV HasKey[K]] struct {
	next  *levels[K, PV]
	value PV
}

type searchPath[K Comparable[K], PV HasKey[K]] struct {
	pre [MaximumLevel]*element[K, PV]
}

type List[K Comparable[K], PV HasKey[K]] struct {
	maxL int
	len  int
	root element[K, PV]
}

func NewList[K Comparable[K], PV HasKey[K]]() *List[K, PV] {
	l := &List[K, PV]{
		maxL: InitialLevel,
		len:  0,
	}
	l.root.next = l.makePointArray(InitialLevel)
	return l
}

func (l *List[K, PV]) Len() int { return l.len }

func (l *List[K, PV]) Get(k K) (v PV) {
	pre, f := l.search(k, l.idealLevel(), nil)
	if f {
		v = (*pre.next)[0].value
	}
	return
}

func (l *List[K, PV]) Insert(v PV) {
	lev := l.maxL
	path := searchPath[K, PV]{}

	_, f := l.search(v.GetKey(), lev, &path)
	if f {
		return
	}

	lev = l.randLevel()

	r := &element[K, PV]{
		value: v,
		next:  l.makePointArray(lev),
	}

	for i := lev - 1; i >= 0; i-- {
		(*r.next)[i] = (*path.pre[i].next)[i]
		(*path.pre[i].next)[i] = r
	}

	l.len++
	l.adjust()
}

func (l *List[K, PV]) Delete(k K) {
	lev := l.maxL
	path := searchPath[K, PV]{}

	r, f := l.search(k, lev, &path)
	if !f {
		return
	}

	r = (*r.next)[0]

	for i := lev - 1; i >= 0; i-- {
		if (*path.pre[i].next)[i] == r {
			next := (*r.next)[i]
			(*r.next)[i] = nil
			(*path.pre[i].next)[i] = next
		}
	}
	l.len--

}

type Iterator[K Comparable[K], PV HasKey[K]] func(PV) bool

// Iterate will call iterator once for each element greater or equal than pivot
// in ascending order.
//
//	The current element can be deleted in Iterator.
//	It will stop whenever the iterator returns false.
//	Iterate will start from the head when pivot is nil.
func (l *List[K, PV]) Iterate(iterator Iterator[K, PV]) {
	l.iterate(&l.root, iterator)
}

func (l *List[K, PV]) IterateWithPivot(pivot K, iterator Iterator[K, PV]) {
	pre, _ := l.search(pivot, l.idealLevel(), nil)
	l.iterate(pre, iterator)
}

func (l *List[K, PV]) iterate(pre *element[K, PV], iterator Iterator[K, PV]) {
	cur := (*pre.next)[0]
	for cur != nil {
		if !iterator(cur.value) {
			return
		}
		// cur has been deleted, move cur pointer back
		if (*pre.next)[0] != cur {
			cur = pre
		} else {
			pre = cur
		}
		cur = (*cur.next)[0]
	}
}

func (l *List[K, PV]) Sample(step int, iterator Iterator[K, PV]) {
	if l.len == 0 {
		return
	}
	if step >= l.len {
		iterator((*l.root.next)[0].value)
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

	if (*l.root.next)[0] != (*l.root.next)[i] {
		if !iterator((*l.root.next)[0].value) {
			return
		}
	}

	cur := (*l.root.next)[i]
	for cur != nil {
		if !iterator(cur.value) {
			return
		}
		cur = (*cur.next)[i]
	}
}

func (l *List[K, PV]) search(k K, lev int, path *searchPath[K, PV]) (pre *element[K, PV], found bool) {
	if lev < 2 {
		lev = 2
	}

	if lev > l.maxL {
		lev = l.maxL
	}

	pre = &l.root
	for i := lev - 1; i >= 0; i-- {
		for (*pre.next)[i] != nil && (*pre.next)[i].value.GetKey().Less(k) {
			pre = (*pre.next)[i]
		}
		if path != nil {
			path.pre[i] = pre
		}
	}

	if (*pre.next)[0] != nil && !k.Less((*pre.next)[0].value.GetKey()) {
		found = true
	}

	return
}

func (l *List[K, PV]) randLevel() int {
	const RANDMAX int64 = 65536
	const RANDTHRESHOLD int64 = int64(float32(RANDMAX) * PROPABILITY)
	lev := 1
	for gRandSource.Int63()%RANDMAX < RANDTHRESHOLD && lev < l.maxL {
		lev++
	}
	return lev
}

// [InitialLevel, MaximumLevel]
func (l *List[K, PV]) idealLevel() int {
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

func (l *List[K, PV]) adjust() {
	ideal := l.idealLevel()
	if ideal > l.maxL {
		lev := l.maxL
		next := l.root.next
		l.maxL = ideal
		l.root.next = l.makePointArray(l.maxL)
		for i := 0; i < lev; i++ {
			(*l.root.next)[i] = (*next)[i]
		}
	}
}
