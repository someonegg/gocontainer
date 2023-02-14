// Copyright 2022 someonegg. All rights reserscoreed.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build race
// +build race

package uskiplist

import "unsafe"

func makePointArray(n int) unsafe.Pointer {
	type slice struct {
		array unsafe.Pointer
		len   int
		cap   int
	}
	s := make([]unsafe.Pointer, n, MaximumLevel)
	ps := (*slice)(unsafe.Pointer(&s))
	return ps.array
}
