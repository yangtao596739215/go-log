package main

import (
	"fmt"

	"github.com/yangtao596739215/go-log/logger"
	"github.com/yangtao596739215/go-log/writer"
)

var GlobalWriterList []writer.LogWriter

var GlobalLogger logger.Log

func init() {
	l := logger.NewLogger("Defult")
	writer, err := writer.NewChanBufferedFileWriter(0, "./log", "logger")
	if err != nil {
		fmt.Println(err)
		return
	}
	l.RegisteWriter(writer)
	GlobalLogger = l
}

func SetLogger(l logger.Log) {
	GlobalLogger = l
}
