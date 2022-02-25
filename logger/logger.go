package logger

import (
	"fmt"
	"runtime"
	"time"

	"github.com/yangtao596739215/go-log/writer"
)

const DEFAULT_LOGGER_NAME = "LOGGER"

//1、定义日志级别
const (
	DEBUG = iota
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

//2、定义log接口
type Log interface {
	Debug(msg string)
	TRACE(msg string)
	INFO(msg string)
	WARNING(msg string)
	ERROR(msg string)
	FATAL(msg string)
}

// 按照日期存储日志信息
type Logger struct {
	level      int8   // 打印日志级别
	name       string // 文件存储路径
	writerList []writer.LogWriter
}

// 构造方法
func NewLogger(name string) *Logger {
	l := &Logger{name: name}
	if name == "" {
		l.name = DEFAULT_LOGGER_NAME
	}
	return l
}

func (l *Logger) RegisteWriter(writer writer.LogWriter) {
	l.writerList = append(l.writerList, writer)
}

// 设置日志级别
func (l *Logger) SetLevel(level int8) {
	l.level = level
}

// 实现Log接口中的方法
func (l *Logger) Debug(msg string) {
	l.writeToLog(DEBUG, msg)
}
func (l *Logger) TRACE(msg string) {
	l.writeToLog(TRACE, msg)
}
func (l *Logger) INFO(msg string) {
	l.writeToLog(INFO, msg)
}
func (l *Logger) WARNING(msg string) {
	l.writeToLog(WARNING, msg)
}
func (l *Logger) ERROR(msg string) {
	l.writeToLog(ERROR, msg)
}
func (l *Logger) FATAL(msg string) {
	l.writeToLog(FATAL, msg)
}

// 向日志中追加内容
func (logger *Logger) writeToLog(level int8, msg string) {
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
		//调栈向上两层，输出调打印日志的方式的文件和行号
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			wrapperdMsg = fmt.Sprintln("[", time.Now().Format("2006-01-02 15:04:05"), "]", "[ERROR]", "[", file, ":", line, "]", "runtime.Caller() fail")
		} else {
			wrapperdMsg = fmt.Sprintln("[", time.Now().Format("2006-01-02 15:04:05"), "]", "[", l, "]", "[", file, ":", line, "]", msg)
		}
	}
	// 遍历所有注册的写入器
	for _, writer := range logger.writerList {
		// 将日志输出到每一个写入器中
		writer.Write(wrapperdMsg)
	}
}
