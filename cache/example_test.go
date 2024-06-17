// Copyright 2022 someonegg. All rights reserscoreed.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cache_test

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/someonegg/gocontainer/cache"
)

type Label struct {
	Type int
	Name string
	Addr string
}

func (l Label) Sprint() string {
	return fmt.Sprintf("%d:%s", l.Type, l.Name)
}

var labelStringerCache = cache.NewStringerCache[Label]()

func (l Label) String() string {
	return labelStringerCache.Get(l)
}

type LabelX Label

func (l LabelX) Sprint() string {
	return fmt.Sprintf("X:%d:%s:%s", l.Type, l.Name, l.Addr)
}

var labelXStringerCache = cache.NewStringerCache[LabelX]()

func (l LabelX) String() string {
	return labelXStringerCache.Get(l)
}

func Example() {
	var s1, s2 string

	s1 = Label{1, "n1", "a1"}.String()
	s2 = Label{1, "n1", "a1"}.String()
	fmt.Println(s1, s2)
	fmt.Println((*reflect.StringHeader)(unsafe.Pointer(&s1)).Data == (*reflect.StringHeader)(unsafe.Pointer(&s2)).Data,
		(*reflect.StringHeader)(unsafe.Pointer(&s1)).Len == (*reflect.StringHeader)(unsafe.Pointer(&s2)).Len)

	s1 = Label{2, "n2", "a2"}.String()
	s2 = Label{2, "n2", "a2"}.String()
	fmt.Println(s1, s2)
	fmt.Println((*reflect.StringHeader)(unsafe.Pointer(&s1)).Data == (*reflect.StringHeader)(unsafe.Pointer(&s2)).Data,
		(*reflect.StringHeader)(unsafe.Pointer(&s1)).Len == (*reflect.StringHeader)(unsafe.Pointer(&s2)).Len)

	fmt.Println()

	s1 = LabelX{1, "n1", "a1"}.String()
	s2 = LabelX{1, "n1", "a1"}.String()
	fmt.Println(s1, s2)
	fmt.Println((*reflect.StringHeader)(unsafe.Pointer(&s1)).Data == (*reflect.StringHeader)(unsafe.Pointer(&s2)).Data,
		(*reflect.StringHeader)(unsafe.Pointer(&s1)).Len == (*reflect.StringHeader)(unsafe.Pointer(&s2)).Len)

	s1 = LabelX{2, "n2", "a2"}.String()
	s2 = LabelX{2, "n2", "a2"}.String()
	fmt.Println(s1, s2)
	fmt.Println((*reflect.StringHeader)(unsafe.Pointer(&s1)).Data == (*reflect.StringHeader)(unsafe.Pointer(&s2)).Data,
		(*reflect.StringHeader)(unsafe.Pointer(&s1)).Len == (*reflect.StringHeader)(unsafe.Pointer(&s2)).Len)

	fmt.Println()

	labelXStringerCache.Clear()

	s2 = LabelX{2, "n2", "a2"}.String()
	fmt.Println(s1, s2)
	fmt.Println((*reflect.StringHeader)(unsafe.Pointer(&s1)).Data == (*reflect.StringHeader)(unsafe.Pointer(&s2)).Data,
		(*reflect.StringHeader)(unsafe.Pointer(&s1)).Len == (*reflect.StringHeader)(unsafe.Pointer(&s2)).Len)

	fmt.Println()

	// Output:
	// 1:n1 1:n1
	// true true
	// 2:n2 2:n2
	// true true
	//
	// X:1:n1:a1 X:1:n1:a1
	// true true
	// X:2:n2:a2 X:2:n2:a2
	// true true
	//
	// X:2:n2:a2 X:2:n2:a2
	// false true
	//
}
