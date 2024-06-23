package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logFile, err := os.OpenFile("./log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("打开日志文件异常")
	}
	defer logFile.Close()
	//log.SetOutput(logFile)

	// 创建一个 MultiWriter，将日志输出到控制台和文件
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	// 设置日志输出到 multiWriter
	log.SetOutput(multiWriter)
	fmt.Println("初始化 Log 成功")
}
