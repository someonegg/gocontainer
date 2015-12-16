// Copyright 2015 someonegg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package list implements a doubly linked list.
// The list's method is similar to standard "container/list".
//
// The difference with standard is:
//
// List management structure(Element) is expect to be
// embedded into user object, not separately alloc.
package list

// Objects implementing the Element interface can be
// inserted into a list.
type Element interface {
	// Wrapped method, for users.
	Next() Element
	Prev() Element

	// Raw method, for list package.
	GetNext() Element
	SetNext(e Element)
	GetPrev() Element
	SetPrev(e Element)
	GetList() *List
	SetList(list *List)
}

// DefElem implements the Element interface.
type DefElem struct {
	next Element
	prev Element
	list *List
}

func (d *DefElem) Next() Element {
	if p := d.next; d.list != nil && p != &d.list.root {
		return p
	}
	return nil
}

func (d *DefElem) Prev() Element {
	if p := d.prev; d.list != nil && p != &d.list.root {
		return p
	}
	return nil
}

func (d *DefElem) GetNext() Element {
	return d.next
}

func (d *DefElem) SetNext(e Element) {
	d.next = e
}

func (d *DefElem) GetPrev() Element {
	return d.prev
}

func (d *DefElem) SetPrev(e Element) {
	d.prev = e
}

func (d *DefElem) GetList() *List {
	return d.list
}

func (d *DefElem) SetList(list *List) {
	d.list = list
}

type List struct {
	root DefElem
	len  int
}

func (l *List) Init() *List {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

func New() *List { return new(List).Init() }

func (l *List) Len() int { return l.len }

func (l *List) Front() Element {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

func (l *List) Back() Element {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

func (l *List) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

func (l *List) insert(e, at Element) Element {
	n := at.GetNext()
	at.SetNext(e)
	e.SetPrev(at)
	e.SetNext(n)
	n.SetPrev(e)
	e.SetList(l)
	l.len++
	return e
}

func (l *List) remove(e Element) Element {
	e.GetPrev().SetNext(e.GetNext())
	e.GetNext().SetPrev(e.GetPrev())
	e.SetNext(nil)
	e.SetPrev(nil)
	e.SetList(nil)
	l.len--
	return e
}

func (l *List) Remove(e Element) Element {
	if e.GetList() == l {
		l.remove(e)
	}
	return e
}

func (l *List) PushFront(e Element) Element {
	l.lazyInit()
	return l.insert(e, &l.root)
}

func (l *List) PushBack(e Element) Element {
	l.lazyInit()
	return l.insert(e, l.root.prev)
}

func (l *List) InsertBefore(e Element, mark Element) Element {
	if mark.GetList() != l {
		return nil
	}
	return l.insert(e, mark.GetPrev())
}

func (l *List) InsertAfter(e Element, mark Element) Element {
	if mark.GetList() != l {
		return nil
	}
	return l.insert(e, mark)
}

func (l *List) MoveToFront(e Element) {
	if e.GetList() != l || l.root.next == e {
		return
	}
	l.insert(l.remove(e), &l.root)
}

func (l *List) MoveToBack(e Element) {
	if e.GetList() != l || l.root.prev == e {
		return
	}
	l.insert(l.remove(e), l.root.prev)
}

func (l *List) MoveBefore(e, mark Element) {
	if e.GetList() != l || e == mark || mark.GetList() != l {
		return
	}
	l.insert(l.remove(e), mark.GetPrev())
}

func (l *List) MoveAfter(e, mark Element) {
	if e.GetList() != l || e == mark || mark.GetList() != l {
		return
	}
	l.insert(l.remove(e), mark)
}

// MergeBackList merge an other list at the back of list l.
// The other list will be reset.
func (l *List) MergeBackList(other *List) {
	l.lazyInit()

	if l == other {
		return
	}

	for {
		e := other.Front()
		if e == nil {
			break
		}
		other.remove(e)
		l.insert(e, l.root.prev)
	}
}

// MergeFrontList merge an other list at the front of list l.
// The other list will be reset.
func (l *List) MergeFrontList(other *List) {
	l.lazyInit()

	if l == other {
		return
	}

	for {
		e := other.Back()
		if e == nil {
			break
		}
		other.remove(e)
		l.insert(e, &l.root)
	}
}
