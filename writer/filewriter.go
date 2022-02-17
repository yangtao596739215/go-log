package writer

import (
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/yangtao596739215/go-log/namerotation"
)

//实现文件的writer

type fileWriter struct {
	writer           io.Writer //文件指针
	mu               sync.RWMutex
	rotationName     string
	ratationMode     int
	pathName         string
	fileName         string
	rotationCallback func()
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
	//如果path指定了一个已经存在的目录，MkdirAll不做任何操作并返回nil。(读r权限值为4,写权限w值为2,执行权限x值为1)
	err := os.MkdirAll(w.pathName, 0777)
	if err != nil {
		return err
	}
	// 创建新的文件 以当前年月日命名.创建文件时要求文件目录必须已经存在
	filewriter, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	fmt.Println(newPath)
	if err != nil {
		return err
	}
	//注册回调关闭fd,close一个已经关闭的fd只会返回err，不会panic，所以此处不需要保护
	w.rotationCallback = func() {
		filewriter.Close()
	}
	w.mu.Lock()
	w.writer = filewriter
	w.mu.Unlock()
	return nil
}

func (w *fileWriter) RotationFileName() {
	if w.needRotetion() {
		w.rotationCallback()
		w.UpdateFile()
	}
}

func (w *fileWriter) needRotetion() bool {
	RotationName, _ := namerotation.GetRotationName(w.ratationMode)
	return w.rotationName == RotationName
}
