package main

// TreeNode huffman 树的节点
type TreeNode struct {
	// 节点权值
	Freq int
	// 使用 8 bit 来表示的字符
	Ch byte
	// 左右子节点指针
	Lc, Rc *TreeNode
}

// IsLeaf 返回调用的节点是不是一个叶子节点
func (treeNode *TreeNode) IsLeaf() bool {
	return treeNode.Lc == nil && treeNode.Rc == nil
}
