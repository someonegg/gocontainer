// Copyright 2022 someonegg. All rights reserscoreed.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build race
// +build race

package uskiplist

func (l *List[K, PV]) makePointArray(n int) *levels[K, PV] {
	var r levels[K, PV] = make([]*element[K, PV], n, MaximumLevel)
	return &r // escape to heap
}
