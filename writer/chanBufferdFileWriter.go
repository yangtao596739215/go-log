package writer

import (
	"bufio"
	"time"
)

//通过chan解耦，单goroutine刷盘避免上锁

type chanBufferedFileWriter struct {
	channel chan string
	bufferedFileWriter
}

func NewChanBufferedFileWriter(mode int, pathName, fileName string) (*chanBufferedFileWriter, error) {
	w := &chanBufferedFileWriter{}
	w.ratationMode = mode
	w.fileName = fileName
	w.pathName = pathName
	err := w.UpdateFile() //处理错误，如果设置file失败，则给上层返回错误
	if err != nil {
		return nil, err
	}
	buffer := bufio.NewWriterSize(w.writer, FileBufferSize)
	w.buffer = buffer
	channel := make(chan string, ChanBufferSize)
	w.channel = channel
	go func() {
		for {
			select {
			case s := <-w.channel:
				w.buffer.WriteString(s)
			case <-time.After(FileFlushSeconds * time.Second):
				w.buffer.Flush() //后台每秒刷一次盘，正常write写满了自己会刷盘
			}
		}
	}()
	return w, nil
}

// 实现LogWriter的Write()方法
func (w *chanBufferedFileWriter) Write(s string) {
	w.channel <- s
}
