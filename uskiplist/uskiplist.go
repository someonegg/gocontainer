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

type levels[K Comparable[K], V HasKey[K]] [MaximumLevel]*element[K, V]

type element[K Comparable[K], V HasKey[K]] struct {
	next  *levels[K, V]
	value *V
}

type searchPath[K Comparable[K], V HasKey[K]] struct {
	pre levels[K, V]
}

type List[K Comparable[K], V HasKey[K]] struct {
	maxL int
	len  int
	root element[K, V]
}

func NewList[K Comparable[K], V HasKey[K]]() *List[K, V] {
	l := &List[K, V]{
		maxL: InitialLevel,
		len:  0,
	}
	l.root.next = (*levels[K, V])(makePointArray(l.maxL))
	return l
}

func (l *List[K, V]) Len() int { return l.len }

func (l *List[K, V]) Get(k K) (v *V) {
	pre, f := l.search(k, l.idealLevel(), nil)
	if f {
		v = pre.next[0].value
	}
	return
}

func (l *List[K, V]) Insert(v *V) {
	path := searchPath[K, V]{}
	lev := l.maxL

	_, f := l.search((*v).GetKey(), lev, &path)
	if f {
		return
	}

	lev = l.randLevel()

	r := &element[K, V]{
		value: v,
		next:  (*levels[K, V])(makePointArray(lev)),
	}

	for i := lev - 1; i >= 0; i-- {
		r.next[i] = path.pre[i].next[i]
		path.pre[i].next[i] = r
	}

	l.len++
	l.adjust()
}

func (l *List[K, V]) Delete(k K) {
	path := searchPath[K, V]{}
	lev := l.maxL

	r, f := l.search(k, lev, &path)
	if !f {
		return
	}

	r = r.next[0]

	for i := lev - 1; i >= 0; i-- {
		if path.pre[i].next[i] == r {
			next := r.next[i]
			r.next[i] = nil
			path.pre[i].next[i] = next
		}
	}
	l.len--

}

type Iterator[K Comparable[K], V HasKey[K]] func(*V) bool

// Iterate will call iterator once for each element greater or equal than pivot
// in ascending order.
//
//	The current element can be deleted in Iterator.
//	It will stop whenever the iterator returns false.
//	Iterate will start from the head when pivot is nil.
func (l *List[K, V]) Iterate(iterator Iterator[K, V]) {
	l.iterate(&l.root, iterator)
}

func (l *List[K, V]) IterateWithPivot(pivot K, iterator Iterator[K, V]) {
	pre, _ := l.search(pivot, l.idealLevel(), nil)
	l.iterate(pre, iterator)
}

func (l *List[K, V]) iterate(pre *element[K, V], iterator Iterator[K, V]) {
	cur := pre.next[0]
	for cur != nil {
		if !iterator(cur.value) {
			return
		}
		// cur has been deleted, move cur pointer back
		if pre.next[0] != cur {
			cur = pre
		} else {
			pre = cur
		}
		cur = cur.next[0]
	}
}

func (l *List[K, V]) Sample(step int, iterator Iterator[K, V]) {
	if l.len == 0 {
		return
	}
	if step >= l.len {
		iterator(l.root.next[0].value)
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
		if !iterator(l.root.next[0].value) {
			return
		}
	}

	cur := l.root.next[i]
	for cur != nil {
		if !iterator(cur.value) {
			return
		}
		cur = cur.next[i]
	}
}

func (l *List[K, V]) search(k K, lev int, path *searchPath[K, V]) (pre *element[K, V], found bool) {
	if lev < 2 {
		lev = 2
	}

	if lev > l.maxL {
		lev = l.maxL
	}

	pre = &l.root
	for i := lev - 1; i >= 0; i-- {
		for pre.next[i] != nil && (*pre.next[i].value).GetKey().Less(k) {
			pre = pre.next[i]
		}
		if path != nil {
			path.pre[i] = pre
		}
	}

	if pre.next[0] != nil && !k.Less((*pre.next[0].value).GetKey()) {
		found = true
	}

	return
}

func (l *List[K, V]) randLevel() int {
	const RANDMAX int64 = 65536
	const RANDTHRESHOLD int64 = int64(float32(RANDMAX) * PROPABILITY)
	lev := 1
	for gRandSource.Int63()%RANDMAX < RANDTHRESHOLD && lev < l.maxL {
		lev++
	}
	return lev
}

// [InitialLevel, MaximumLevel]
func (l *List[K, V]) idealLevel() int {
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

func (l *List[K, V]) adjust() {
	ideal := l.idealLevel()
	if ideal > l.maxL {
		lev := l.maxL
		next := l.root.next
		l.maxL = ideal
		l.root.next = (*levels[K, V])(makePointArray(l.maxL))
		for i := 0; i < lev; i++ {
			l.root.next[i] = next[i]
		}
	}
}
