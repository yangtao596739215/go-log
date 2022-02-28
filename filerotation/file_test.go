package filerotation

import (
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
	"time"
)

//测试在大量并发写的情况下，替换底层的文件
func Test_write(t *testing.T) {

	f, err := os.Create("text.log")
	if err != nil {
		fmt.Println(err.Error())
		t.Error("xxx")
	}

	var w io.Writer
	w = f

	wg := &sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			//每个循环写1000行
			for j := 0; j < 1000; j++ {
				w.Write([]byte(fmt.Sprintf("xxxxx第%d个\n", i)))
			}
			wg.Done()
		}(i)
	}

	f1, err := os.Create("text1.log")
	if err != nil {
		fmt.Println(err.Error())
		t.Error("xxx")
	}
	time.Sleep(100 * time.Millisecond) //file1写100ms
	w = f1
	wg.Wait() //等待写完，统计行数

	//成功：1.没有数据丢失 2.没有报错
}
