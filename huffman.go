package main

import (
	"flag"
	"strings"
)

const R = 256

var (
	method   = flag.String("m", "", "m 方法 \n encode 对目标文件进行编码 \n decode 对目标文件进行解码 \n check 对目标文件进行 md5 检查，比对文件是否一致")
	filePath = flag.String("f", "", "f 文件名,文件名（第二个文件名仅在检查时可用）")
)

func main() {
	flag.Parse()
	filePaths := strings.Split(*filePath, ",")
	length := len(filePaths)
	if length == 1 {
		switch *method {
		case "encode":
			EncodeDoor(filePaths[0])
		case "decode":
			DecodeDoor(filePaths[0])
		default:
			flag.PrintDefaults()
		}
	} else if length == 2 {
		if *method == "check" {
			Check(filePaths[0], filePaths[1])
		} else {
			flag.PrintDefaults()
		}
	} else {
		flag.PrintDefaults()
	}
}
