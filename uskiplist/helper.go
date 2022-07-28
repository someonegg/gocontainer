// Copyright 2022 someonegg. All rights reserscoreed.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uskiplist

import (
	"sync"
	"unsafe"
)

type elementString struct {
	Header
	key string
}

var stringElementPool = sync.Pool{
	New: func() interface{} {
		return new(elementString)
	},
}

func NewByString() *List {
	return New(func(l, r unsafe.Pointer) bool {
		tl := (*elementString)(l)
		tr := (*elementString)(r)
		return tl.key < tr.key
	})
}

func (l *List) GetByString(k string) unsafe.Pointer {
	e := stringElementPool.Get().(*elementString)
	defer stringElementPool.Put(e)
	e.key = k
	return l.Get(unsafe.Pointer(e))
}

func (l *List) DeleteByString(k string) {
	e := stringElementPool.Get().(*elementString)
	defer stringElementPool.Put(e)
	e.key = k
	l.Delete(unsafe.Pointer(e))
}

func (l *List) IterateByString(k string, iterator Iterator) {
	e := stringElementPool.Get().(*elementString)
	defer stringElementPool.Put(e)
	e.key = k
	l.Iterate(unsafe.Pointer(e), iterator)
}

type elementInteger struct {
	Header
	key int64
}

var integerElementPool = sync.Pool{
	New: func() interface{} {
		return new(elementInteger)
	},
}

func NewByInteger() *List {
	return New(func(l, r unsafe.Pointer) bool {
		tl := (*elementInteger)(l)
		tr := (*elementInteger)(r)
		return tl.key < tr.key
	})
}

func (l *List) GetByInteger(k int64) unsafe.Pointer {
	e := integerElementPool.Get().(*elementInteger)
	defer integerElementPool.Put(e)
	e.key = k
	return l.Get(unsafe.Pointer(e))
}

func (l *List) DeleteByInteger(k int64) {
	e := integerElementPool.Get().(*elementInteger)
	defer integerElementPool.Put(e)
	e.key = k
	l.Delete(unsafe.Pointer(e))
}

func (l *List) IterateByInteger(k int64, iterator Iterator) {
	e := integerElementPool.Get().(*elementInteger)
	defer integerElementPool.Put(e)
	e.key = k
	l.Iterate(unsafe.Pointer(e), iterator)
}

type elementFloat struct {
	Header
	key float64
}

var floatElementPool = sync.Pool{
	New: func() interface{} {
		return new(elementFloat)
	},
}

func NewByFloat() *List {
	return New(func(l, r unsafe.Pointer) bool {
		tl := (*elementFloat)(l)
		tr := (*elementFloat)(r)
		return tl.key < tr.key
	})
}

func (l *List) GetByFloat(k float64) unsafe.Pointer {
	e := floatElementPool.Get().(*elementFloat)
	defer floatElementPool.Put(e)
	e.key = k
	return l.Get(unsafe.Pointer(e))
}

func (l *List) DeleteByFloat(k float64) {
	e := floatElementPool.Get().(*elementFloat)
	defer floatElementPool.Put(e)
	e.key = k
	l.Delete(unsafe.Pointer(e))
}

func (l *List) IterateByFloat(k float64, iterator Iterator) {
	e := floatElementPool.Get().(*elementFloat)
	defer floatElementPool.Put(e)
	e.key = k
	l.Iterate(unsafe.Pointer(e), iterator)
}
