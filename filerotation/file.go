package filerotation

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type FileManager interface {
	GetFile() *os.File
}

type fileManager struct {
	curFile     *os.File
	fileName    string
	filePath    string
	rotationMod int
	saveTime    time.Duration
}

//对外暴露方法，获取底层保存的文件,之后更新的时候进行指针替换
func (m *fileManager) GetFile() *os.File {
	return m.curFile
}

type FileSaveTime time.Duration

const (
	FileSaveOneHour  = 1 * time.Hour
	FileSaveOneDay   = 24 * time.Hour
	FileSaveThreeDay = 72 * time.Hour
)

func NewFileManager(filename, filepath string, rotationmod int, savetime time.Duration) (FileManager, error) {
	manager := &fileManager{
		fileName:    filename,
		filePath:    filepath,
		rotationMod: rotationmod,
		saveTime:    savetime,
	}
	err := manager.initFile()
	if err != nil {
		return nil, err
	}
	go manager.deleteExpiredFile()
	go manager.updateFile()
	return manager, nil
}

func (m *fileManager) initFile() error {
	RotationName, _ := GetRotationName(m.rotationMod)
	fileName := m.fileName + RotationName
	newPath := path.Join(m.filePath, fileName)
	//如果path指定了一个已经存在的目录，MkdirAll不做任何操作并返回nil。(读r权限值为4,写权限w值为2,执行权限x值为1)
	err := os.MkdirAll(m.filePath, 0777)
	if err != nil {
		return err
	}
	// 创建新的文件 以当前年月日命名.创建文件时要求文件目录必须已经存在
	tmpFile, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	//替换文件
	m.curFile = tmpFile
	return nil
}

//读取对应目录下的文件，删除过期的文件
func (m *fileManager) deleteExpiredFile() {
	for {
		//sleep放后面，每次启动的时候，先把过期的文件删除了
		deletedName, err := GetRotationNameByTime(time.Now().Add(-m.saveTime), m.rotationMod)
		if err != nil {
			m.curFile.WriteString("getrotationName failed:" + err.Error())
			return
		}

		rd, err := ioutil.ReadDir(m.filePath)
		if err != nil {
			m.curFile.WriteString("read dir failed:" + err.Error())
			return
		}

		for _, fi := range rd {
			if !fi.IsDir() {
				fmt.Println(fi.Name())
				//按照文件名的字符串大小比较
				if fi.Name() < deletedName {
					os.Remove(fi.Name())
				}
			}
		}

		//每小时执行一次，去删除一天前的文件
		time.Sleep(1 * time.Hour)
	}
}

//根据循环模式，每天或每小时执行一次
func (m *fileManager) updateFile() {
	for {
		m.initFile()
		now := time.Now()
		// 计算下一次执行的时间(整hour或整day)
		var next time.Time
		switch m.rotationMod {
		case ROTATION_BY_DAY:
			next = now.Add(24 * time.Hour)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		case ROTATION_BY_HOUR:
			next := now.Add(time.Hour)
			next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())
		default:
			m.curFile.WriteString("updateFile faied \n")
		}
		t := time.NewTimer(next.Sub(now))
		<-t.C
	}
}
