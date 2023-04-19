package main

import (
	"container/heap"
	"testing"
	"time"
)

func TestPriorityQueue_Pop(t *testing.T) {
	now := time.Now()
	getNow = func() time.Time {
		return now.Add(10 * time.Minute)
	}
	// Create an array of time items
	ti := [6]*TimeItem{}
	for i := 0; i < 6; i++ {
		ti[i] = &TimeItem{
			value: PlayerUpdate{
				WhoAmI: i,
			},
			priority: now.Add(time.Duration(i) * time.Second),
		}
	}
	/*
			items := []*TimeItem{
			{value: PlayerUpdate{}, priority: now.Add(5 * time.Second)},
			{value: PlayerUpdate{}, priority: now.Add(1 * time.Second)},
			{value: PlayerUpdate{}, priority: now.Add(10 * time.Second)},
			{value: PlayerUpdate{}, priority: now.Add(3 * time.Second)},
		}
	*/
	pq := make(PriorityQueue, 0)

	heap.Init(&pq)
	for _, val := range ti {
		heap.Push(&pq, val)
	}

	// Pop the items off the priority queue and check that they are in the expected order
	expected := []int{0, 1, 2, 3, 4, 5}
	for _, i := range expected {
		item := heap.Pop(&pq).(*TimeItem)
		if item == nil {
			t.Errorf("item is nil, expected non-nil")
		}
		if item.value.WhoAmI != i {
			t.Errorf("item value whoami is %d, expected %d", item.value.WhoAmI, i)
		}
	}

	// Make sure there are no more items on the priority queue
	if pq.Len() != 0 {
		t.Errorf("pq len is %d, expected 0", pq.Len())
	}
}
