package writer

import (
	"io"
	"sync"
	"time"

	"github.com/yangtao596739215/go-log/filerotation"
)

//实现文件的writer
type fileWriter struct {
	writer      io.Writer  //文件指针
	mu          sync.Mutex //并发写锁，防止数据错误
	filemanager filerotation.FileManager
}

// 实现LogWriter的Write()方法
func (w *fileWriter) Write(s string) {
	w.mu.Lock()
	w.filemanager.GetFile().WriteString(s)
	w.mu.Unlock()
}

func NewFileWriter(mode int, pathName, fileName string) (*fileWriter, error) {
	w := &fileWriter{}
	m, err := filerotation.NewFileManager(fileName, pathName, mode, 72*time.Hour)
	if err != nil {
		return nil, err
	}
	w.filemanager = m
	w.writer = m.GetFile()
	return w, nil
}
