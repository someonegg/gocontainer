// Copyright 2015 someonegg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package queue contains some queue implements.
package queue

import (
	"container/list"
	"github.com/someonegg/gox/syncx"
	"sync"
)

// Queue is a double-ended FIFO list.
//
// You can init the queue manually, see Init method.
type Queue struct {
	list   list.List
	locker sync.Locker
}

// Init the queue manually, with a locker (can be nil).
func (q *Queue) Init(l sync.Locker) {
	q.list.Init()
	q.locker = l
}

func (q *Queue) Len() int {
	return q.list.Len()
}

func (q *Queue) PushFront(e interface{}) {
	if q.locker != nil {
		q.locker.Lock()
		defer q.locker.Unlock()
	}
	q.list.PushFront(e)
}

func (q *Queue) PushBack(e interface{}) {
	if q.locker != nil {
		q.locker.Lock()
		defer q.locker.Unlock()
	}
	q.list.PushBack(e)
}

func (q *Queue) PopFront() interface{} {
	if q.locker != nil {
		q.locker.Lock()
		defer q.locker.Unlock()
	}
	e := q.list.Front()
	if e == nil {
		return nil
	}
	return q.list.Remove(e)
}

func (q *Queue) PopBack() interface{} {
	if q.locker != nil {
		q.locker.Lock()
		defer q.locker.Unlock()
	}
	e := q.list.Back()
	if e == nil {
		return nil
	}
	return q.list.Remove(e)
}

func (q *Queue) Lock() {
	if q.locker != nil {
		q.locker.Lock()
	}
}

func (q *Queue) Unlock() {
	if q.locker != nil {
		q.locker.Unlock()
	}
}

func (q *Queue) ObtainList() *list.List {
	return &q.list
}

// EventQueue is a queue with event(a read-only chan), the
// event will return data if the list is not empty.
//
// You must init the queue manually, see Init method.
type EventQueue struct {
	Queue
	event syncx.Event
}

// Init the queue manually, with a locker (can be nil).
func (q *EventQueue) Init(l sync.Locker) {
	q.Queue.Init(l)
	q.event = syncx.NewEvent()
}

func (q *EventQueue) PushFront(e interface{}) {
	q.Queue.PushFront(e)
	q.SetEvent()
}

func (q *EventQueue) PushBack(e interface{}) {
	q.Queue.PushBack(e)
	q.SetEvent()
}

func (q *EventQueue) PopFront() interface{} {
	e := q.Queue.PopFront()
	if q.Queue.Len() > 0 {
		q.SetEvent()
	}
	return e
}

func (q *EventQueue) PopBack() interface{} {
	e := q.Queue.PopBack()
	if q.Queue.Len() > 0 {
		q.SetEvent()
	}
	return e
}

func (q *EventQueue) Event() syncx.EventR {
	return q.event.R()
}

func (q *EventQueue) SetEvent() {
	q.event.Set()
}
