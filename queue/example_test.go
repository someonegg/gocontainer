// Copyright 2015 someonegg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue_test

import (
	"fmt"
	"github.com/someonegg/gocontainer/queue"
)

func Example() {
	// An empty queue and put some numbers in it.
	var q queue.EventQueue
	q.Init(nil)
	q.PushFront(1)
	q.PushFront(2)
	q.PushFront(3)
	q.PushFront(4)

	select {
	case <-q.Event():
	}

	for e := q.PopFront(); e != nil; e = q.PopFront() {
		fmt.Println(e.(int))
	}

	// Output:
	// 4
	// 3
	// 2
	// 1
}
