package logger

import (
	"bytes"
	"errors"
	"sync"
)

var ErrCreateLog = errors.New("create buf error")

//定义全局的pool实现日志的临时存储
var bufPool = &sync.Pool{
	New: func() interface{} {
		//函数体内容
		return bytes.NewBuffer([]byte{})
	},
}

// 按照日期存储日志信息
type BufferLogger struct {
	level  int8 // 打印日志级别
	buffer *bytes.Buffer
}

func NewBufferdLogger(level int8) (*BufferLogger, error) {
	b, ok := bufPool.Get().(*bytes.Buffer)
	if !ok {
		return nil, ErrCreateLog
	}
	bLog := &BufferLogger{
		buffer: b,
	}
	return bLog, nil
}

// 设置日志级别
func (l *BufferLogger) SetLevel(level int8) {
	l.level = level
}

// 实现Log接口中的方法
func (l *BufferLogger) Debug(msg string) {
	l.writeToBuffer(DEBUG, msg)
}
func (l *BufferLogger) TRACE(msg string) {
	l.writeToBuffer(TRACE, msg)
}
func (l *BufferLogger) INFO(msg string) {
	l.writeToBuffer(INFO, msg)
}
func (l *BufferLogger) WARNING(msg string) {
	l.writeToBuffer(WARNING, msg)
}
func (l *BufferLogger) ERROR(msg string) {
	l.writeToBuffer(ERROR, msg)
}
func (l *BufferLogger) FATAL(msg string) {
	l.writeToBuffer(FATAL, msg)
}

func (logger *BufferLogger) writeToBuffer(level int8, msg string) {
	// 当前级别大于日志级别才输出 否则不输出
	if level >= logger.level {
		logger.buffer.WriteString(msg)
	}
}

func (l *BufferLogger) Flush() string {
	defer func() {
		l.buffer.Reset() //清空字符串，不会释放底层[]byte
		bufPool.Put(l.buffer)
		l.buffer = nil
	}()
	return l.buffer.String()
}
