// Copyright 2026 someonegg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sortedmap_test

import (
	"reflect"
	"testing"

	"github.com/someonegg/gocontainer/sortedmap"
)

func TestOrderedMap(t *testing.T) {
	m := sortedmap.NewOrdered[string, int]()

	if got := m.Get("missing"); got != nil {
		t.Fatalf("Get missing = %v, want nil", *got)
	}

	a := m.Set("b", 2)
	m.Set("a", 1)
	m.Set("c", 3)

	if m.Len() != 3 {
		t.Fatalf("Len = %d, want 3", m.Len())
	}

	a2 := m.Set("b", 20)
	if a2 != a {
		t.Fatal("Set existing key changed value address")
	}
	if *m.Get("b") != 20 {
		t.Fatalf("Get b = %d, want 20", *m.Get("b"))
	}

	var keys []string
	var values []int
	m.Range(func(k string, v *int) bool {
		keys = append(keys, k)
		values = append(values, *v)
		*v *= 10
		return true
	})

	if !reflect.DeepEqual(keys, []string{"a", "b", "c"}) {
		t.Fatalf("Range keys = %v, want [a b c]", keys)
	}
	if !reflect.DeepEqual(values, []int{1, 20, 3}) {
		t.Fatalf("Range values = %v, want [1 20 3]", values)
	}
	if *m.Get("a") != 10 || *m.Get("b") != 200 || *m.Get("c") != 30 {
		t.Fatalf("Range value mutation did not persist")
	}

	keys = keys[:0]
	values = values[:0]
	m.RangeFrom("b", func(k string, v *int) bool {
		keys = append(keys, k)
		values = append(values, *v)
		return true
	})
	if !reflect.DeepEqual(keys, []string{"b", "c"}) {
		t.Fatalf("RangeFrom keys = %v, want [b c]", keys)
	}
	if !reflect.DeepEqual(values, []int{200, 30}) {
		t.Fatalf("RangeFrom values = %v, want [200 30]", values)
	}

	old, ok := m.Delete("b")
	if !ok || old != 200 {
		t.Fatalf("Delete b = (%d, %v), want (200, true)", old, ok)
	}
	if got := m.Get("b"); got != nil {
		t.Fatalf("Get deleted b = %v, want nil", *got)
	}
	if old, ok := m.Delete("b"); ok || old != 0 {
		t.Fatalf("Delete missing b = (%d, %v), want (0, false)", old, ok)
	}

	m.Clear()
	if m.Len() != 0 {
		t.Fatalf("Len after Clear = %d, want 0", m.Len())
	}
	if got := m.Get("a"); got != nil {
		t.Fatalf("Get after Clear = %v, want nil", *got)
	}
}

func TestOrderedMapDeleteCurrentDuringRange(t *testing.T) {
	m := sortedmap.NewOrdered[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)
	m.Set("d", 4)

	var keys []string
	m.Range(func(k string, _ *int) bool {
		keys = append(keys, k)
		if k == "b" || k == "c" {
			if _, ok := m.Delete(k); !ok {
				t.Fatalf("Delete current key %q failed", k)
			}
		}
		return true
	})

	if !reflect.DeepEqual(keys, []string{"a", "b", "c", "d"}) {
		t.Fatalf("Range keys = %v, want [a b c d]", keys)
	}
	if m.Len() != 2 || m.Get("b") != nil || m.Get("c") != nil {
		t.Fatalf("current entries were not deleted correctly")
	}
	if got := *m.Get("a"); got != 1 {
		t.Fatalf("Get a = %d, want 1", got)
	}
	if got := *m.Get("d"); got != 4 {
		t.Fatalf("Get d = %d, want 4", got)
	}
}

type testKey struct {
	major int
	minor int
}

func (k testKey) Less(k2 testKey) bool {
	if k.major != k2.major {
		return k.major < k2.major
	}
	return k.minor < k2.minor
}

func TestMap(t *testing.T) {
	m := sortedmap.New[testKey, string]()

	m.Set(testKey{major: 2, minor: 1}, "c")
	m.Set(testKey{major: 1, minor: 2}, "b")
	m.Set(testKey{major: 1, minor: 1}, "a")

	p := m.Get(testKey{major: 1, minor: 2})
	if p == nil || *p != "b" {
		t.Fatalf("Get custom key = %v, want b", p)
	}

	p2 := m.Set(testKey{major: 1, minor: 2}, "bb")
	if p2 != p {
		t.Fatal("Set existing custom key changed value address")
	}

	var values []string
	m.Range(func(_ testKey, v *string) bool {
		values = append(values, *v)
		return true
	})
	if !reflect.DeepEqual(values, []string{"a", "bb", "c"}) {
		t.Fatalf("Range custom key values = %v, want [a bb c]", values)
	}

	values = values[:0]
	m.RangeFrom(testKey{major: 1, minor: 2}, func(_ testKey, v *string) bool {
		values = append(values, *v)
		return false
	})
	if !reflect.DeepEqual(values, []string{"bb"}) {
		t.Fatalf("RangeFrom stop values = %v, want [bb]", values)
	}
}

func TestMapDeleteCurrentDuringRangeFrom(t *testing.T) {
	m := sortedmap.New[testKey, string]()
	a := testKey{major: 1, minor: 1}
	b := testKey{major: 1, minor: 2}
	c := testKey{major: 2, minor: 1}
	d := testKey{major: 3, minor: 1}

	m.Set(a, "a")
	m.Set(b, "b")
	m.Set(c, "c")
	m.Set(d, "d")

	var values []string
	m.RangeFrom(b, func(k testKey, v *string) bool {
		values = append(values, *v)
		if k == b || k == c {
			if _, ok := m.Delete(k); !ok {
				t.Fatalf("Delete current key %+v failed", k)
			}
		}
		return true
	})

	if !reflect.DeepEqual(values, []string{"b", "c", "d"}) {
		t.Fatalf("RangeFrom values = %v, want [b c d]", values)
	}
	if m.Len() != 2 || m.Get(b) != nil || m.Get(c) != nil {
		t.Fatalf("current entries were not deleted correctly")
	}
	if got := *m.Get(a); got != "a" {
		t.Fatalf("Get a = %q, want a", got)
	}
	if got := *m.Get(d); got != "d" {
		t.Fatalf("Get d = %q, want d", got)
	}
}
