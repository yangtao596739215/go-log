package writer

import (
	"bufio"
	"time"

	"github.com/yangtao596739215/go-log/filerotation"
)

//设置文件写缓冲，减少刷盘次数。
type bufferedFileWriter struct {
	buffer *bufio.Writer
	fileWriter
}

func NewBufferdFileWriter(mode int, pathName, fileName string) (*bufferedFileWriter, error) {
	w := &bufferedFileWriter{}
	m, err := filerotation.NewFileManager(fileName, pathName, mode, 72*time.Hour)
	if err != nil {
		return nil, err
	}
	buffer := bufio.NewWriterSize(m.GetFile(), FileBufferSize)
	w.buffer = buffer
	w.filemanager = m
	go w.flushBuffer()
	return w, nil
}

//每秒将buffer进行一次刷盘
func (b *bufferedFileWriter) flushBuffer() {
	for {
		time.Sleep(BufferFlushSeconds * time.Second)
		b.mu.Lock()
		b.buffer.Flush() //后台每秒刷一次盘，正常write写满了自己会刷盘
		b.mu.Unlock()
	}
}

func (b *bufferedFileWriter) Write(s string) {
	b.mu.Lock()
	b.buffer.WriteString(s)
	b.mu.Unlock()
}
