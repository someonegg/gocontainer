// Copyright 2026 someonegg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sortedmap implements generic sorted maps based on uskiplist.
package sortedmap

import (
	"github.com/someonegg/gocontainer/cmp"
	"github.com/someonegg/gocontainer/uskiplist"
)

type entry[K any, V any] struct {
	uskiplist.Embedder[entry[K, V]]

	key   K
	value V
}

func (e *entry[K, V]) Key() K {
	return e.key
}

// Map is a sorted map whose keys define their ordering with Less.
//
// The address of a value is stable until the entry is deleted or the map is
// cleared. Keys are not addressable through this API.
type Map[K cmp.Key[K], V any] struct {
	list uskiplist.List[K, entry[K, V], *entry[K, V]]
}

// New creates and initializes a new sorted map.
func New[K cmp.Key[K], V any]() *Map[K, V] {
	m := &Map[K, V]{}
	m.Init()
	return m
}

// Init initializes the map.
func (m *Map[K, V]) Init() {
	m.list.Init()
}

// Len returns number of entries in the map.
func (m *Map[K, V]) Len() int {
	return m.list.Len()
}

// Clear removes all entries from the map.
func (m *Map[K, V]) Clear() {
	m.Init()
}

// Get returns a pointer to the value associated with k.
//
// It returns nil when k is not found. The returned pointer remains valid until
// the entry is deleted or the map is cleared.
func (m *Map[K, V]) Get(k K) *V {
	e := m.list.Get(k)
	if e == nil {
		return nil
	}
	return &e.value
}

// Set sets the value for k and returns a pointer to the stored value.
//
// When k already exists, Set overwrites the existing value without changing its
// address.
func (m *Map[K, V]) Set(k K, v V) *V {
	e := m.list.Get(k)
	if e != nil {
		e.value = v
		return &e.value
	}

	e = &entry[K, V]{
		key:   k,
		value: v,
	}
	m.list.Insert(e)
	return &e.value
}

// Delete removes k from the map and returns the removed value.
func (m *Map[K, V]) Delete(k K) (old V, ok bool) {
	e := m.list.Get(k)
	if e == nil {
		return old, false
	}

	old = e.value
	m.list.Delete(k)
	return old, true
}

// Range calls fn once for each entry in ascending key order.
//
// It stops whenever fn returns false. The current entry may be deleted during
// iteration, and values may be changed through the provided pointer.
func (m *Map[K, V]) Range(fn func(K, *V) bool) {
	m.list.Iterate(func(e *entry[K, V]) bool {
		return fn(e.key, &e.value)
	})
}

// RangeFrom calls fn once for each entry with key greater than or equal to
// pivot, in ascending key order.
//
// It stops whenever fn returns false. The current entry may be deleted during
// iteration, and values may be changed through the provided pointer.
func (m *Map[K, V]) RangeFrom(pivot K, fn func(K, *V) bool) {
	m.list.IterateFrom(pivot, func(e *entry[K, V]) bool {
		return fn(e.key, &e.value)
	})
}

// OrderedMap is a sorted map whose keys are ordered with the < operator.
//
// The address of a value is stable until the entry is deleted or the map is
// cleared. Keys are not addressable through this API.
type OrderedMap[K cmp.Ordered, V any] struct {
	list uskiplist.ListO[K, entry[K, V], *entry[K, V]]
}

// NewOrdered creates and initializes a new sorted map for ordered keys.
func NewOrdered[K cmp.Ordered, V any]() *OrderedMap[K, V] {
	m := &OrderedMap[K, V]{}
	m.Init()
	return m
}

// Init initializes the map.
func (m *OrderedMap[K, V]) Init() {
	m.list.Init()
}

// Len returns number of entries in the map.
func (m *OrderedMap[K, V]) Len() int {
	return m.list.Len()
}

// Clear removes all entries from the map.
func (m *OrderedMap[K, V]) Clear() {
	m.Init()
}

// Get returns a pointer to the value associated with k.
//
// It returns nil when k is not found. The returned pointer remains valid until
// the entry is deleted or the map is cleared.
func (m *OrderedMap[K, V]) Get(k K) *V {
	e := m.list.Get(k)
	if e == nil {
		return nil
	}
	return &e.value
}

// Set sets the value for k and returns a pointer to the stored value.
//
// When k already exists, Set overwrites the existing value without changing its
// address.
func (m *OrderedMap[K, V]) Set(k K, v V) *V {
	e := m.list.Get(k)
	if e != nil {
		e.value = v
		return &e.value
	}

	e = &entry[K, V]{
		key:   k,
		value: v,
	}
	m.list.Insert(e)
	return &e.value
}

// Delete removes k from the map and returns the removed value.
func (m *OrderedMap[K, V]) Delete(k K) (old V, ok bool) {
	e := m.list.Get(k)
	if e == nil {
		return old, false
	}

	old = e.value
	m.list.Delete(k)
	return old, true
}

// Range calls fn once for each entry in ascending key order.
//
// It stops whenever fn returns false. The current entry may be deleted during
// iteration, and values may be changed through the provided pointer.
func (m *OrderedMap[K, V]) Range(fn func(K, *V) bool) {
	m.list.Iterate(func(e *entry[K, V]) bool {
		return fn(e.key, &e.value)
	})
}

// RangeFrom calls fn once for each entry with key greater than or equal to
// pivot, in ascending key order.
//
// It stops whenever fn returns false. The current entry may be deleted during
// iteration, and values may be changed through the provided pointer.
func (m *OrderedMap[K, V]) RangeFrom(pivot K, fn func(K, *V) bool) {
	m.list.IterateFrom(pivot, func(e *entry[K, V]) bool {
		return fn(e.key, &e.value)
	})
}
