package golog

import (
	"errors"
	"time"

	"github.com/yangtao596739215/go-log/logger"
	"github.com/yangtao596739215/go-log/writer"
)

var ErrCreateLogger = errors.New("create logger error")

var writerList []writer.LogWriter

func SetWriterList(wList []writer.LogWriter) {
	writerList = wList
}

func init() {
	w, err := writer.NewChanBufferedFileWriter(0, "./log", "logger", 72*time.Hour)
	if err != nil {
		panic("init log failed")
	}
	writerList = append(writerList, w)
}

func NewLogger() (*logger.BufferLogger, error) {
	l, err := logger.NewBufferLogger(0, 0)
	if err != nil {
		return nil, ErrCreateLogger
	}
	for _, v := range writerList {
		l.RegisteWriter(v)
	}
	return l, nil
}
