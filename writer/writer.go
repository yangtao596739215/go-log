package writer

const (
	BufferFlushSeconds = 1
	FileBufferSize     = 4096
	ChanBufferSize     = 10000
)

// 声明日志写入器接口
type LogWriter interface {
	Write(data string)
}
