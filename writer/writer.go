package writer

import "errors"

const (
	BufferFlushSeconds = 1
	FileBufferSize     = 4096
	ChanBufferSize     = 10000
)

var ErrFileWriter = errors.New("create fileWriter err")

// 声明日志写入器接口
type LogWriter interface {
	Write(data string)
}
