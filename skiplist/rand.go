// Copyright 2022 someonegg. All rights reserscoreed.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package skiplist

import "math/rand"

// A SplitMix64 provides the SplitMix64 algorithm and implements
// math/rand.Source64. Can be seeded to any value.
// https://pkg.go.dev/nullprogram.com/x/rng#SplitMix64
type splitMix64 uint64

var _ rand.Source64 = (*splitMix64)(nil)

func (s *splitMix64) Seed(seed int64) {
	*s = splitMix64(seed)
}

func (s *splitMix64) Uint64() uint64 {
	*s += 0x9e3779b97f4a7c15
	z := uint64(*s)
	z ^= z >> 30
	z *= 0xbf58476d1ce4e5b9
	z ^= z >> 27
	z *= 0x94d049bb133111eb
	z ^= z >> 31
	return z
}

func (s *splitMix64) Int63() int64 {
	return int64(s.Uint64() >> 1)
}
