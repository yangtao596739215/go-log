package writer

import (
	"bufio"
	"time"
)

//设置文件写缓冲，减少刷盘次数。

type bufferedFileWriter struct {
	buffer *bufio.Writer
	fileWriter
}

func NewBufferdFileWriter(mode int, pathName, fileName string) (*bufferedFileWriter, error) {
	w := &bufferedFileWriter{}
	w.ratationMode = mode
	w.fileName = fileName
	w.pathName = pathName
	err := w.UpdateFile() //处理错误，如果设置file失败，则给上层返回错误
	if err != nil {
		return nil, err
	}
	buffer := bufio.NewWriterSize(w.writer, FileBufferSize)
	w.buffer = buffer
	go func() {
		time.Sleep(FileFlushSeconds * time.Second)
		w.mu.Lock()
		w.buffer.Flush() //后台每秒刷一次盘，正常write写满了自己会刷盘
		w.mu.Unlock()
	}()
	return w, nil
}
