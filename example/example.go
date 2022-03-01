package main

import (
	"fmt"
	"time"

	golog "github.com/yangtao596739215/go-log"
	"github.com/yangtao596739215/go-log/writer"
)

func main() {
	//使用默认的filewriter
	l, err := golog.NewLogger()
	if err != nil {
		fmt.Println("err")

	}
	l.INFO("hahaha") //此处可以把logger传入业务方法中使用
	l.Flush()
	time.Sleep(3 * time.Second)

	//自定义filewriter
	cbfw, err := writer.NewChanBufferedFileWriter(1, "./log", "custom", 72*time.Hour)
	if err != nil {
		fmt.Println("err")
	}
	golog.SetWriterList([]writer.LogWriter{cbfw})
	l2, err := golog.NewLogger()
	if err != nil {
		fmt.Println("err")

	}
	l2.INFO("l2l2l2l2") //此处可以把logger传入业务方法中使用
	l2.Flush()
	time.Sleep(3 * time.Second)
}
