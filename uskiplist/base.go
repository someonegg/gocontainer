// Copyright 2022 someonegg. All rights reserscoreed.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package uskiplist implements a generic skiplist using
// unusual operations to minimize memory and references.
package uskiplist

import (
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

type level1[E any] [1]*E
type leveln[E any] [MaximumLevel]*E

// Embedder should be embedded in the skiplist element struct.
type Embedder[E any] struct {
	next unsafe.Pointer
}

func (e *Embedder[E]) ptNext() *unsafe.Pointer {
	return &e.next
}

func (e *Embedder[E]) l1Next() *level1[E] {
	return (*level1[E])(unsafe.Pointer(&e.next))
}

func (e *Embedder[E]) lnNext() *leveln[E] {
	return (*leveln[E])(e.next)
}

type Iterator[E any] func(*E) bool

type searchPath[E any] [MaximumLevel]**E

type listBase[E any] struct {
	maxL int
	len  int
	root *leveln[E]
	rnd  splitMix64
}

func (l *listBase[E]) init() {
	l.maxL = InitialLevel
	l.len = 0
	l.root = (*leveln[E])(makePointArray(InitialLevel))
	l.rnd = splitMix64(time.Now().Unix())
}

// [InitialLevel, MaximumLevel]
func (l *listBase[E]) idealLevel() int {
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

func (l *listBase[E]) adjust() {
	ideal := l.idealLevel()
	if ideal > l.maxL {
		lev := l.maxL
		root := l.root
		l.maxL = ideal
		l.root = (*leveln[E])(makePointArray(l.maxL))
		for i := 0; i < lev; i++ {
			l.root[i] = root[i]
		}
	}
}

func (l *listBase[E]) randLevel() int {
	const RANDMAX int64 = 65536
	const RANDTHRESHOLD int64 = int64(float32(RANDMAX) * PROPABILITY)
	lev := 1
	for l.rnd.Int63()%RANDMAX < RANDTHRESHOLD && lev < l.maxL {
		lev++
	}
	return lev
}
