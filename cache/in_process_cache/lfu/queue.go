package lfu

import "tcpserver/cache/in_process_cache"

type entry struct {
	key string

	val string

	weight int

	index int
}

func (e *entry) Len() int {
	return in_process_cache.CalcLen(e.val) + 4 + 4
}

type queue []*entry

func (q queue) Len() int {
	return len(q)
}

func (q queue) Less(i, j int) bool {
	return q[i].weight < q[j].weight
}

func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[i].index = j
}

func (q *queue) Push(x interface{}) {
	n := len(*q)
	en := x.(*entry)
	en.index = n
	*q = append(*q, en)
}

func (q queue) Pop() interface{} {
	panic("implement me")
}
