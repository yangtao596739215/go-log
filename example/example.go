package main

import (
	golog "github.com/yangtao596739215/go-log"
	"github.com/yangtao596739215/go-log/logger"
)

func main() {

	//case1  中间不收到输入换行的话，只打印一行日志。日志先放入buffer，然后再通过chan写文件。日志中的文件信息只有当前这个函数的
	bl, err := logger.NewBufferdLogger(3) //每个请求一个对象
	if err != nil {
		panic("logger init err")
	}

	bl.WARNING("xxx")

	golog.GlobalLogger.INFO(bl.Flush())

	//case2 一次打印一行，每行都带有文件信息
	golog.GlobalLogger.INFO("xxxx") //全局一个对象，底层是chan，可以并发写

}
