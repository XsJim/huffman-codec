package main

import (
	"container/heap"
	"fmt"
	"io"
	"log"
	"time"
)

//1. 获取读取器
//2. 生成 byte 统计数组
//3. 获取优先队列-------
//4. 非 0 统计元素组合成节点并入队
//5. 生成树
//6. 生成写入器
//7. 通过 树写入器（树，写入器） 将树写入文件
//8. 将原文件byte长度写入到输出文件-----防止后续解码文件时读出补足位 bit
//9. 通过 获取编码表（树）（编码表） 获得编码表
//10. 读取器游标置 0
//11. 通过 编码器（编码表，写入器，读取器） 写入文件

// EncodeDoor 压缩入口
func EncodeDoor(filePath string) {
	// 输出
	fmt.Println("-> 正在进行压缩...")
	begin := time.Now()
	Encode(filePath)
	end := time.Now()
	use := end.Sub(begin)
	// 输出
	fmt.Printf("-> 压缩完成，用时：%s\n", use)
	// 计算压缩率
	fmt.Println("-> 正在计算压缩率...")
	c := contrast(filePath, filePath+".huff")

	// 输出
	fmt.Printf("-> 压缩率为：%.2f%%\n", c*100)
}

// contrast 计算出文件 2 和文件 1 的大小向除的结果并返回
func contrast(filePath, encodeFilePath string) float64 {
	// 打开第一个文件
	reader1, err := NewReader(filePath)
	if err != nil {
		log.Fatal(err)
	}
	// 打开第二个文件
	reader2, err := NewReader(encodeFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// 返回文件压缩率
	return float64(reader2.BufferLen()) / float64(reader1.BufferLen())
}

func Encode(filePath string) {
	// 获取读取器
	reader, err := NewReader(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// 生成 byte 统计数组
	count := make([]int, R)
	for {
		b, err := reader.Byte()
		count[b]++
		if err == io.EOF {
			break
		}
	}

	// 获取优先队列
	pq := &PriorityQueue{}

	// 统计数量非 0 的字符组装成节点入队
	for i, v := range count {
		if v != 0 {
			heap.Push(pq, &TreeNode{v, byte(i), nil, nil})
		}
	}

	// 生成树
	for pq.Len() > 1 {
		lc := heap.Pop(pq).(*TreeNode)
		rc := heap.Pop(pq).(*TreeNode)
		heap.Push(pq, &TreeNode{lc.Freq + rc.Freq, 0, lc, rc})
	}
	treeRoot := heap.Pop(pq).(*TreeNode)

	// 生成写入器
	// 默认的输出文件是当前文件路径加上 .huff
	writer, err := NewWriter(filePath + ".huff")
	if err != nil {
		log.Fatal(err)
	}

	// 通过 树写入器（树，写入器） 将树写入文件
	treeWriter(treeRoot, writer)

	//将原文件byte长度写入到输出文件-----防止后续解码文件时读出补足位 bit
	writer.Int(reader.BufferLen() / 8)

	// 通过 获取编码表（树）（编码表） 获得编码表
	codeScheme := getCodeScheme(treeRoot)

	// 读取器游标置 0
	reader.Reset()

	// 通过 编码器（编码表，写入器，读取器） 写入文件
	encode(codeScheme, writer, reader)

	// 刷出缓冲区内容并关闭文件
	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

// treeWriter 树写入器，将树结构输出到写入器中
// 注意：huffman 树，要么节点是叶子节点，要么节点一定有两个孩子
// 这是根据构造规则得出的结论
func treeWriter(treeRoot *TreeNode, writer *Writer) {
	// 以先序遍历顺序遍历树，但输出时还附带了节点左右子树是否存在

	curNode := treeRoot
	stack := &Stack{}

	for curNode != nil || stack.Size() > 0 {
		for curNode != nil {
			// 第一次遇到节点
			if curNode.IsLeaf() {
				// 如果是树叶，则写入 1 并且将其表示的字符写入
				writer.Bit(true)
				writer.Byte(curNode.Ch)
			} else {
				// 否则写入 0
				writer.Bit(false)
			}

			stack.Push(curNode)
			curNode = curNode.Lc
		}

		if stack.Size() > 0 {
			curNode = stack.Pop().(*TreeNode).Rc
		}
	}
}

// getCodeScheme 该方法遍历哈夫曼树，生成一张编码表并返回
func getCodeScheme(treeRoot *TreeNode) (codeScheme [][]bool) {
	codeScheme = make([][]bool, R)
	buildCodeScheme(treeRoot, codeScheme, make([]bool, 0))

	return
}

// buildCodeScheme 递归填充编码表
// 先不考虑文本只有同一字符的情况
func buildCodeScheme(treeNode *TreeNode, codeScheme [][]bool, curCode []bool) {
	if treeNode.IsLeaf() {
		codeScheme[treeNode.Ch] = make([]bool, len(curCode))
		copy(codeScheme[treeNode.Ch], curCode)
	} else {
		curCode = append(curCode, false)
		buildCodeScheme(treeNode.Lc, codeScheme, curCode)
		curCode = append(curCode[:len(curCode)-1], true)
		buildCodeScheme(treeNode.Rc, codeScheme, curCode)
	}
}

// encode 该方法读取字节，匹配编码，并将编码输出到写入器
func encode(codeScheme [][]bool, writer *Writer, reader *Reader) {
	for {
		b, err := reader.Byte()
		if err == io.EOF {
			break
		}
		writer.Bits(codeScheme[b])
	}
}
