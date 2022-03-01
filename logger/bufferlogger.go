package logger

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/yangtao596739215/go-log/writer"
)

// 定义日志级别
const (
	DEBUG = iota
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

// 定义log接口
type Log interface {
	Debug(format string, args ...interface{})
	TRACE(format string, args ...interface{})
	INFO(format string, args ...interface{})
	WARNING(format string, args ...interface{})
	ERROR(format string, args ...interface{})
	FATAL(format string, args ...interface{})
}

var ErrCreateLog = errors.New("create buf error")

// 定义全局的pool实现日志的临时存储
var bufPool = &sync.Pool{
	New: func() interface{} {
		return &BufferLogger{
			buffer: bytes.NewBuffer([]byte{}),
		}
	},
}

// 一个请求创建一个对象，flush方法自动回收
type BufferLogger struct {
	level      int32 // 打印日志级别
	logid      int32
	buffer     *bytes.Buffer
	mu         sync.Mutex
	writerList []writer.LogWriter
}

func NewBufferLogger(logid, level int32) (*BufferLogger, error) {
	b, ok := bufPool.Get().(*BufferLogger)
	if !ok {
		return nil, ErrCreateLog
	}
	return b, nil
}

// 设置日志id
func (l *BufferLogger) SetLogid(logid int32) {
	l.logid = logid
}

// 设置日志级别
func (l *BufferLogger) SetLevel(level int32) {
	l.level = level
}

// 实现Log接口中的方法
func (l *BufferLogger) Debug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.writeToBuffer(DEBUG, msg)
}
func (l *BufferLogger) TRACE(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.writeToBuffer(TRACE, msg)
}
func (l *BufferLogger) INFO(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.writeToBuffer(INFO, msg)
}
func (l *BufferLogger) WARNING(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.writeToBuffer(WARNING, msg)
}
func (l *BufferLogger) ERROR(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.writeToBuffer(ERROR, msg)
}
func (l *BufferLogger) FATAL(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.writeToBuffer(FATAL, msg)
}

func (logger *BufferLogger) writeToBuffer(level int32, msg string) {
	var l, wrapperdMsg string
	switch level {
	case DEBUG:
		l = "DEBUG"
	case TRACE:
		l = "TRACE"
	case INFO:
		l = "INFO"
	case WARNING:
		l = "WARNING"
	case ERROR:
		l = "ERROR"
	case FATAL:
		l = "FATAL"
	}
	// 当前级别大于日志级别才输出 否则不输出
	if level >= logger.level {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			wrapperdMsg = fmt.Sprintln("[", time.Now().Format("2006-01-02 15:04:05"), "]", "[ERROR]", "[", file, ":", line, "]", "[", logger.logid, "]", "runtime.Caller() fail")
		} else {
			wrapperdMsg = fmt.Sprintln("[", time.Now().Format("2006-01-02 15:04:05"), "]", "[", l, "]", "[", file, ":", line, "]", "[", logger.logid, "]", msg)
		}
		//避免多线程写入造成数据错乱
		logger.mu.Lock()
		logger.buffer.WriteString(wrapperdMsg)
		logger.mu.Unlock()
	}
}

// 刷到下一层，对象复用
func (l *BufferLogger) Flush() {
	defer func() {
		l.buffer.Reset() //清空字符串，不会释放底层[]byte
		bufPool.Put(l)   //把自己放回池子
	}()
	l.writeToLog(l.buffer.String())
}

func (l *BufferLogger) RegisteWriter(writer writer.LogWriter) {
	l.writerList = append(l.writerList, writer)
}

// 向日志中追加内容
func (logger *BufferLogger) writeToLog(msg string) {
	// 遍历所有注册的写入器
	for _, writer := range logger.writerList {
		// 将日志输出到每一个写入器中
		writer.Write(msg)
	}
}
