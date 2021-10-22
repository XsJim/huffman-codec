package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

//1. 获取读取器
//2. 通过 解析树（读取器）树根 来获得树根
//3. 获取写入器
//4. 通过 解码器（树根，读取器，写入器） 将文件从读取器读出并解码到写入器中

// DecodeDoor 解压入口
func DecodeDoor(filePath string) {
	fmt.Println("-> 正在解压中，请稍后...")
	begin := time.Now()
	Decode(filePath)
	end := time.Now()

	use := end.Sub(begin)

	fmt.Printf("-> 解压完成，用时：%s\n", use)
}

// Decode 解码传入的文件
func Decode(filePath string) {
	// 获取读取器
	reader, err := NewReader(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// 通过 解析树（读取器）树根 来获得树根
	treeRoot := decodeTree(reader)

	// 获取写入器
	writer, err := NewWriter("huffman-out-" + filePath[:strings.LastIndex(filePath, ".huff")])

	// 通过 解码器（树根，读取器，写入器） 将文件从读取器读出并解码到写入器中
	decode(treeRoot, reader, writer)

	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

// decodeTree 从读取器内读入数据，组合成一颗 huffman 编码树
// 并返回这个树的树根
func decodeTree(reader *Reader) (treeRoot *TreeNode) {
	treeRoot = readTreeNode(reader)
	return
}

// readTreeNode 该方法顺序读取内容并递归的构造树
func readTreeNode(reader *Reader) (treeNode *TreeNode) {
	// 首先读入 1 bit
	bit, err := reader.Bit()
	if err != nil {
		log.Fatal(err)
	}
	if bit {
		// 这是一个叶节点，接下来 8 位是它所表示的字符
		b, err := reader.Byte()
		if err != nil {
			log.Fatal(err)
		}
		treeNode = &TreeNode{Ch: b}
	} else {
		// 这只是一个路径
		// 继续读取树，生成左右子树
		treeNode = &TreeNode{
			Lc: readTreeNode(reader),
			Rc: readTreeNode(reader),
		}
	}

	return treeNode
}

// decode 对读取器中顺序向后的内容进行解码
// 并将解码内容输入到写入器中
func decode(treeRoot *TreeNode, reader *Reader, writer *Writer) {
	// 首先读出未压缩文件的 byte 长度
	byteLen, err := reader.Int()

	if err != nil {
		log.Fatal("decode ", err)
	}

	// 已经输出的 byte 的计数器
	curByteLen := 0

	curNode := treeRoot
	for {
		bit, err := reader.Bit()
		if err != nil {
			// 解码完毕
			return
		}

		// 根据读入的编码，在树中移动
		if bit {
			curNode = curNode.Rc
		} else {
			curNode = curNode.Lc
		}

		// 如果到达了树叶，说明该段编码已经可以解码
		if curNode.IsLeaf() {
			// 将解码输入写入器
			writer.Byte(curNode.Ch)

			// 计数器加一
			curByteLen++

			// 如果已经达到原文件 byte 数，说明读取器后续的编码都是补足位，无需读入
			if curByteLen == byteLen {
				return
			}

			// 重置当前位置
			curNode = treeRoot
		}
	}
}
