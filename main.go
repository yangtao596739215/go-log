package main

import (
	"fmt"
	"time"

	"github.com/yangtao596739215/go-log/logger"
	"github.com/yangtao596739215/go-log/writer"
)

func main() {

	l := logger.NewLogger("test")
	writer, err := writer.NewChanBufferedFileWriter(0, "./log", "test")
	if err != nil {
		fmt.Println(err)
		return
	}
	l.RegisteWriter(writer)
	l.SetLevel(0)
	l.Debug("xxx")
	time.Sleep(3 * time.Second)
}
