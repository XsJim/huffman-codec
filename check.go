package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
)

// Check 该方法检测两个文件是否相同（通过计算 md5）
// 并且会输出相关的信息
func Check(filePath1, filePath2 string) {
	file1, err := os.Open(filePath1)
	if err != nil {
		log.Fatal(err)
	}
	defer file1.Close()

	file2, err := os.Open(filePath2)
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	h := md5.New()

	_, err = io.Copy(h, file1)
	if err != nil {
		log.Fatal(err)
	}
	file1Md5 := fmt.Sprintf("%x", h.Sum(nil))

	h.Reset()
	_, err = io.Copy(h, file2)
	if err != nil {
		log.Fatal(err)
	}

	file2Md5 := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println("<huffman 文件对比信息>")

	fmt.Print("结果-文件内容是否相同:")

	if file1Md5 == file2Md5 {
		fmt.Println("相同")
	} else {
		fmt.Println("不相同")
	}

	fmt.Println("<huffman 文件信息>")

	fmt.Println("\t文件-1\t\t\t\t\t文件-2")
	fmt.Printf("文件名\t%s\t\t\t\t%s\n", file1.Name(), file2.Name())
	fmt.Printf("md5\t%s\t%s\n", file1Md5, file2Md5)
}
