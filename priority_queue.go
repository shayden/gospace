package main

import "time"

var getNow func() time.Time = func() time.Time {
	return time.Now()
}

type TimeItem struct {
	value    PlayerUpdate
	priority time.Time
	index    int
}

type PriorityQueue []*TimeItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority.Before(pq[j].priority)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*TimeItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	if pq.Len() == 0 {
		return nil
	}
	item := (*pq)[0]
	if getNow().Before(item.priority) {
		return nil
	}
	old := *pq
	n := len(old)
	item = old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

/*
func main() {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	heap.Push(&pq, &TimeItem{value: "item1", priority: time.Now().Add(1 * time.Minute)})
	heap.Push(&pq, &TimeItem{value: "item2", priority: time.Now().Add(5 * time.Minute)})
	heap.Push(&pq, &TimeItem{value: "item3", priority: time.Now().Add(3 * time.Minute)})

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*TimeItem)
		fmt.Printf("Value: %s, Priority: %s\n", item.value, item.priority)
	}
}
*/
