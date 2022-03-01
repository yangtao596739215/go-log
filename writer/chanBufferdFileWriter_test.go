package writer

import (
	"testing"
	"time"
)

func Test_chanBufferedFileWriter_writeAndFlushBuffer(t *testing.T) {
	cbfw, err := NewChanBufferedFileWriter(1, "./log", "test", 72*time.Hour)
	if err != nil {
		t.Error(err)
	}

	cbfw.Write("xxx")
	time.Sleep(10 * time.Second) //通过日志查看flush的次数来判断是否符合预期
}
