package writer

// 声明日志写入器接口
type LogWriter interface {
	Write(data string)
}
