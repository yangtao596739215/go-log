# go-log
go语言的开发的一个基础日志库



#### 通用Logger实现，实现日志级别的区分

#### Logger与Writer解耦，方便自定义writer，将日志送到不同的地方

#### 实现三种writer
1.fileWriter                基本文件的writer，支持文件名按照时间切分
2.bufferedFileWriter        带有buffer的writer，减少磁盘的io次数
3.chanBufferedFileWriter    带有chan的writer，优化了buffer的锁
