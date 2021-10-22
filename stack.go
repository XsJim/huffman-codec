package main

// Stack 该结构体用来表示一个栈，内部只包含一个切片，存放元素的容量和内存有关
// 前些天实现的栈结构
type Stack struct {
	arr []interface{}
}

// Push 栈的添加方法
func (s *Stack) Push(x interface{}) {
	s.arr = append(s.arr, x)
}

// Pop 弹出当前栈顶元素，该方法返回两个值，第一个值是栈顶元素（如果不存在则为 nil），第二个值是栈中是否有元素（弹出是否成功）
func (s *Stack) Pop() (top interface{}) {
	top = s.arr[len(s.arr)-1]
	s.arr = s.arr[:len(s.arr)-1]
	return
}

// Size 调用此方法返回栈内元素数量
func (s *Stack) Size() int {
	return len(s.arr)
}
