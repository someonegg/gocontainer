// Copyright 2022 someonegg. All rights reserscoreed.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build race
// +build race

package uskiplist

import "unsafe"

func makePointArray(n int) unsafe.Pointer {
	slice := make([]unsafe.Pointer, n, MaximumLevel)
	array := unsafe.SliceData(slice)
	return unsafe.Pointer(array)
}
