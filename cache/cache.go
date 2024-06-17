// Copyright 2022 someonegg. All rights reserscoreed.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cache provides several cache utilities.
package cache

import (
	"sync"
)

type Stringer interface {
	comparable

	Sprint() string
}

type StringerCache[O Stringer] struct {
	mu    sync.RWMutex
	cache map[O]string
}

func NewStringerCache[O Stringer]() *StringerCache[O] {
	return &StringerCache[O]{
		cache: make(map[O]string),
	}
}

func (c *StringerCache[O]) Get(o O) string {
	c.mu.RLock()
	s, ok := c.cache[o]
	if ok {
		c.mu.RUnlock()
		return s
	}
	c.mu.RUnlock()

	c.mu.Lock()

	s, ok = c.cache[o]
	if !ok {
		s = o.Sprint()
		c.cache[o] = s
	}

	c.mu.Unlock()
	return s
}

func (c *StringerCache[O]) Clear() {
	c.mu.Lock()
	c.cache = make(map[O]string)
	c.mu.Unlock()
}
