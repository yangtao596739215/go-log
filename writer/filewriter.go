package writer

import (
	"io"
	"time"

	"github.com/yangtao596739215/go-log/filerotation"
)

//实现文件的writer

type fileWriter struct {
	writer io.Writer //文件指针
}

// 实现LogWriter的Write()方法
func (w *fileWriter) Write(s string) {
	filerotation.GetFile().WriteString(s)
}

func NewFileWriter(mode int, pathName, fileName string) error {
	return filerotation.InitFileManager(fileName, pathName, mode, 72*time.Hour)
}
