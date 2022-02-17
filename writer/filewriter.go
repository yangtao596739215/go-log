package writer

import (
	"io"
	"os"
	"path"
	"sync"

	"github.com/yangtao596739215/go-log/namerotation"
)

//实现文件的writer

type fileWriter struct {
	writer       io.Writer //文件指针
	mu           sync.RWMutex
	rotationName string
	ratationMode int
	pathName     string
	fileName     string
}

// 实现LogWriter的Write()方法
func (w *fileWriter) Write(s string) {
	if w.writer == nil {
		return // 日志文件没有准备好，则直接返回
	}
	if w.needRotetion() {
		w.UpdateFile() //这里不处理错误，如果rotate失败，则写老的file
	}
	w.mu.Lock()
	w.writer.Write([]byte(s))
	w.mu.Unlock()
}

func NewFileWriter(mode int, pathName, fileName string) (*fileWriter, error) {
	w := &fileWriter{}
	w.ratationMode = mode
	w.fileName = fileName
	w.pathName = pathName
	err := w.UpdateFile() //处理错误，如果设置file失败，则给上层返回错误
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *fileWriter) UpdateFile() error {
	RotationName, _ := namerotation.GetRotationName(w.ratationMode)
	fileName := w.fileName + RotationName
	newPath := path.Join(w.pathName, fileName)
	// 创建新的文件 以当前年月日命名
	filewriter, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return ErrFileWriter
	}
	w.mu.Lock()
	w.writer = filewriter
	w.mu.Unlock()
	return nil
}

func (w *fileWriter) RotationFileName() {
	if w.needRotetion() {
		w.UpdateFile()
	}
}

func (w *fileWriter) needRotetion() bool {
	RotationName, _ := namerotation.GetRotationName(w.ratationMode)
	return w.rotationName == RotationName
}
