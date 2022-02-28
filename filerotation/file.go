package filerotation

import (
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"
)

type fileManager struct {
	curFile     *os.File
	fileName    string
	filePath    string
	rotationMod int
	saveTime    time.Duration
	mu          sync.RWMutex
}

var manager *fileManager

//对外暴露方法，获取底层保存的文件,之后更新的时候进行指针替换
func GetFile() *os.File {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	return manager.curFile
}

type FileSaveTime time.Duration

const (
	FileSaveOneHour  = 1 * time.Hour
	FileSaveOneDay   = 24 * time.Hour
	FileSaveThreeDay = 72 * time.Hour
)

func InitFileManager(filename, filepath string, rotationmod int, savetime time.Duration) error {
	manager = &fileManager{
		fileName:    filename,
		filePath:    filepath,
		rotationMod: rotationmod,
		saveTime:    savetime,
	}
	err := manager.initFile()
	if err != nil {
		return err
	}
	go manager.deleteExpiredFile()
	go manager.updateFile()
	return nil
}

func (m *fileManager) initFile() error {
	RotationName, _ := GetRotationName(m.rotationMod)
	fileName := m.fileName + RotationName
	newPath := path.Join(m.filePath, fileName)
	//如果path指定了一个已经存在的目录，MkdirAll不做任何操作并返回nil。(读r权限值为4,写权限w值为2,执行权限x值为1)
	err := os.MkdirAll(newPath, 0777)
	if err != nil {
		return err
	}
	// 创建新的文件 以当前年月日命名.创建文件时要求文件目录必须已经存在
	tmpFile, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	//加锁替换文件
	manager.mu.Lock()
	manager.curFile = tmpFile
	manager.mu.Unlock()
	return nil
}

//读取对应目录下的文件，删除过期的文件
func (m *fileManager) deleteExpiredFile() {
	for {
		//每小时执行一次，去删除一天前的文件
		time.Sleep(1 * time.Hour)
		deletedName, err := GetRotationNameByTime(time.Now().Add(-m.saveTime), m.rotationMod)
		if err != nil {
			m.curFile.WriteString("getrotationName failed:" + err.Error())
			return
		}

		rd, err := ioutil.ReadDir(manager.filePath)
		if err != nil {
			m.curFile.WriteString("read dir failed:" + err.Error())
			return
		}

		for _, fi := range rd {
			if !fi.IsDir() {
				//按照文件名的字符串大小比较
				if fi.Name() < deletedName {
					os.Remove(fi.Name())
				}
			}
		}
	}
}

//根据循环模式，每天或每小时执行一次
func (m *fileManager) updateFile() {
	for {
		//todo
		m.initFile()
		//每小时执行一次
		now := time.Now()
		// 计算下一次执行的时间
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
