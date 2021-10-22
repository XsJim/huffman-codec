package main

type PriorityQueue []*TreeNode

// 实现 sort.Interface 接口

func (p PriorityQueue) Len() int {
	return len(p)
}

func (p PriorityQueue) Less(i, j int) bool {
	return p[i].Freq < p[j].Freq
}

func (p PriorityQueue) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// 实现 container.Interface 接口

func (p *PriorityQueue) Push(x interface{}) {
	*p = append(*p, x.(*TreeNode))
}

func (p *PriorityQueue) Pop() interface{} {
	old := *p
	n := p.Len()
	x := old[n-1]
	*p = old[:n-1]
	return x
}
