package main

import (
	"fmt"
	"io"
	"log"
	"testing"
)

func TestReaderAndWrite(t *testing.T) {
	reader, err := NewReader("西遊記.txt")

	if err != nil {
		fmt.Println(err)
	} else {
		b := make([]byte, 1)
		writer, err := NewWriter("out.txt")
		if err != nil {
			log.Fatal(err)
		}
		for {
			b[0], err = reader.Byte()
			if err == io.EOF {
				break
			}
			writer.Byte(b[0])
		}
		err = writer.Flush()
		if err != nil {
			log.Println(err)
		}
	}
}
