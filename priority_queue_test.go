package main

import (
	"container/heap"
	"log"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	pq := &PriorityQueue{}
	node1 := &TreeNode{2, 0, nil, nil}
	node2 := &TreeNode{1, 0, nil, nil}
	node3 := &TreeNode{3, 0, nil, nil}
	node4 := &TreeNode{4, 0, nil, nil}

	heap.Push(pq, node1)
	heap.Push(pq, node2)
	heap.Push(pq, node3)
	heap.Push(pq, node4)

	for pq.Len() > 0 {
		log.Println(heap.Pop(pq).(*TreeNode).Freq)
	}
}
