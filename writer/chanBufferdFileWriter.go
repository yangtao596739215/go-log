package writer

import (
	"bufio"
	"fmt"
	"time"

	"github.com/yangtao596739215/go-log/filerotation"
)

//通过chan解耦，单goroutine刷盘避免上锁
type chanBufferedFileWriter struct {
	channel chan string
	bufferedFileWriter
}

func NewChanBufferedFileWriter(mode int, pathName, fileName string, saveTime time.Duration) (*chanBufferedFileWriter, error) {
	w := &chanBufferedFileWriter{}
	m, err := filerotation.NewFileManager(fileName, pathName, mode, saveTime)
	if err != nil {
		fmt.Println("new manager err:" + err.Error())
		return nil, err
	}
	w.filemanager = m
	buffer := bufio.NewWriterSize(m.GetFile(), FileBufferSize)
	w.buffer = buffer
	channel := make(chan string, ChanBufferSize)
	w.channel = channel
	go w.writeAndFlushBuffer()
	return w, nil
}

// 实现LogWriter的Write()方法
func (w *chanBufferedFileWriter) Write(s string) {
	w.channel <- s
}

// 写或者刷新buffer(在单个goroutine中处理buffer，可以不用加锁)
func (w *chanBufferedFileWriter) writeAndFlushBuffer() {
	//创建定时器，每隔1秒后，定时器就会给channel发送一个事件(当前时间)
	ticker := time.NewTicker(BufferFlushSeconds * time.Second)
	for {
		select {
		case s := <-w.channel:
			w.buffer.WriteString(s)
			ticker.Reset(2 * time.Second) //没有写操作2s后刷盘，可以注释掉。在日志量能快速写满4096的时候开启，可以提高性能
		case <-ticker.C:
			w.buffer.Flush() //后台每秒刷一次盘，正常write写满了自己会刷盘
		}
	}
}
