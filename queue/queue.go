// Copyright 2015 someonegg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package queue contains some queue implements.
//
// According to the design, the element push to queue can not be nil.
package queue

import (
	"errors"
	"github.com/someonegg/gocontainer/list"
	"github.com/someonegg/goutility/chanutil"
	"sync"
)

var ErrNilElement = errors.New("queue's element is nil")

// Queue is a double-ended FIFO list.
type Queue struct {
	list list.List
}

func (q *Queue) List() *list.List {
	return &q.list
}

func (q *Queue) Len() int {
	return q.list.Len()
}

func (q *Queue) PushFront(e list.Element) {
	if e == nil {
		panic(ErrNilElement)
	}
	q.list.PushFront(e)
}

func (q *Queue) PushBack(e list.Element) {
	if e == nil {
		panic(ErrNilElement)
	}
	q.list.PushBack(e)
}

func (q *Queue) PopFront() list.Element {
	e := q.list.Front()
	if e == nil {
		return nil
	}
	return q.list.Remove(e)
}

func (q *Queue) PopBack() list.Element {
	e := q.list.Back()
	if e == nil {
		return nil
	}
	return q.list.Remove(e)
}

// SynQueue is a syn queue.
type SynQueue struct {
	inner Queue
	lock  sync.Mutex
}

func (q *SynQueue) Lock() {
	q.lock.Lock()
}

func (q *SynQueue) List() *list.List {
	return q.inner.List()
}

func (q *SynQueue) Unlock() {
	q.lock.Unlock()
}

func (q *SynQueue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.inner.Len()
}

func (q *SynQueue) PushFront(e list.Element) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.inner.PushFront(e)
}

func (q *SynQueue) PushBack(e list.Element) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.inner.PushBack(e)
}

func (q *SynQueue) PopFront() list.Element {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.inner.PopFront()
}

func (q *SynQueue) PopBack() list.Element {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.inner.PopBack()
}

// EvtQueue is a queue with event(a read-only chan), the
// event will return data if the list is not empty.
//
// Users must call Init() before use.
type EvtQueue struct {
	inner SynQueue
	event chanutil.Event
}

func (q *EvtQueue) Init() {
	q.event = chanutil.NewEvent()
}

func (q *EvtQueue) Lock() {
	q.inner.Lock()
}

func (q *EvtQueue) List() *list.List {
	return q.inner.List()
}

func (q *EvtQueue) Unlock() {
	q.inner.Unlock()
}

func (q *EvtQueue) Len() int {
	return q.inner.Len()
}

func (q *EvtQueue) PushFront(e list.Element) {
	q.inner.PushFront(e)
	q.SetEvent()
}

func (q *EvtQueue) PushBack(e list.Element) {
	q.inner.PushBack(e)
	q.SetEvent()
}

func (q *EvtQueue) PopFront() list.Element {
	e := q.inner.PopFront()
	if q.inner.Len() > 0 {
		q.SetEvent()
	}
	return e
}

func (q *EvtQueue) PopBack() list.Element {
	e := q.inner.PopBack()
	if q.inner.Len() > 0 {
		q.SetEvent()
	}
	return e
}

func (q *EvtQueue) Event() chanutil.EventR {
	return q.event.R()
}

func (q *EvtQueue) SetEvent() {
	q.event.Set()
}
