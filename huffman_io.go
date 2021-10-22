package main

import (
	"io"
	"log"
	"os"
)

// Reader 读取器
type Reader struct {
	// 缓存区，缓存区的长度应该始终为 8 的整数倍
	buffer []bool

	// 读取器游标当前位置，起始位置下标为 0
	cursor int

	// 缓存区长度
	length int
}

// NewReader 构造一个读取器
// 构造器直接读取整个文件
func NewReader(filePath string) (*Reader, error) {
	// 尝试打开文件
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	// 生成读取器
	r := new(Reader)
	// 初始读取器
	r.buffer = make([]bool, 0, 512)
	r.cursor = 0
	r.length = 0
	// 读取文件
	b := make([]byte, 128)
	for {
		n, err := file.Read(b)
		r.length += n
		for i := 0; i < n; i++ {
			for j := 0; j < 8; j++ {
				r.buffer = append(r.buffer, (b[i]&1) == 1)
				b[i] >>= 1
			}
		}

		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
	}

	r.length *= 8
	return r, nil
}

// Bit 该方法返回读取器当前游标处的 bit ，并将游标向后移动
// 如果已经到达缓存区尾部，err = io.EOF
func (reader *Reader) Bit() (bit bool, err error) {
	if reader.cursor == reader.length {
		err = io.EOF
	} else {
		bit = reader.buffer[reader.cursor]
		reader.cursor++
	}

	return
}

// Byte 该方法返回从读取器当前游标开始的向后 8 bit 组成的 byte
// 如果不足 8 位或者已经到达末尾，err = io.EOF
func (reader *Reader) Byte() (b byte, err error) {
	if reader.cursor > reader.length-8 {
		err = io.EOF
	} else {
		for i := 0; i < 8; i++ {
			if reader.buffer[reader.cursor] {
				b |= 1 << i
			}
			reader.cursor++
		}
	}

	return b, err
}

// Int 从读取器中取出一个 32 位 int
func (reader *Reader) Int() (i int, err error) {
	if reader.cursor > reader.length-32 {
		err = io.EOF
	} else {
		for j := 0; j < 4; j++ {
			b, err := reader.Byte()
			if err != nil {
				log.Fatal("reader.Int ", err)
			}
			i |= int(b) << (j * 8)
		}
	}

	return
}

// Reset 重置游标到起始位置
func (reader *Reader) Reset() {
	reader.cursor = 0
}

// BufferLen 返回这个读取器缓存区的长度
func (reader *Reader) BufferLen() int {
	return reader.length
}

// Writer 写入器
type Writer struct {
	// 写入器缓存
	buffer []bool

	// 文件句柄
	file *os.File
}

// NewWriter 构造一个写入器
// 构造器创建并打开文件，获得文件句柄，便于之后直接写入文件
func NewWriter(filePath string) (*Writer, error) {
	writer := new(Writer)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	writer.file = file

	writer.buffer = make([]bool, 0, 512)

	return writer, nil
}

// Bit 该方法将一个 bit 写入缓存
func (writer *Writer) Bit(bit bool) {
	writer.buffer = append(writer.buffer, bit)
}

// Byte 该方法将一个 byte 写入缓存
func (writer *Writer) Byte(b byte) {
	for i := 0; i < 8; i++ {
		writer.Bit((b & 1) == 1)
		b >>= 1
	}
}

// Bits 该方法将一串 bit 写入缓存
func (writer *Writer) Bits(bits []bool) {
	for _, bit := range bits {
		writer.Bit(bit)
	}
}

// Int 将一个整形写入缓存，这个整形将被输出低 32 位
func (writer *Writer) Int(x int) {
	for i := 0; i < 32; i++ {
		writer.Bit((x & 1) == 1)
		x >>= 1
	}
}

// Flush 该方法将写入器的缓存区数据写入到文件当中
// 执行该方法后，文件句柄被关闭
func (writer *Writer) Flush() error {
	length := len(writer.buffer)
	i := 0
	b := make([]byte, 0, 128)
	for {
		if i+8 < length {
			b = append(b, getByte(writer.buffer[i:i+8]))
			i += 8
		} else {
			b = append(b, getByte(writer.buffer[i:]))
			break
		}

		if len(b) == cap(b) {
			_, err := writer.file.Write(b)
			if err != nil {
				return err
			}
			b = b[:0]
		}
	}
	_, err := writer.file.Write(b)
	writer.file.Close()
	return err
}

// getByte 从一个 bit 切片中组合出一个 byte
func getByte(bits []bool) (b byte) {
	for i, bit := range bits {
		if bit {
			b |= 1 << i
		}
	}

	return b
}
