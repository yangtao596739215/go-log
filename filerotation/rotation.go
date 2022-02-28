package filerotation

import (
	"errors"
	"time"
)

var ErrNameRotation = errors.New("get rotation name err")

//名字循环
const (
	ROTATION_BY_DAY  = iota //按天循环
	ROTATION_BY_HOUR        //按小时循环
)

func GetRotationName(mode int) (string, error) {
	now := time.Now()
	switch mode {
	case ROTATION_BY_DAY:
		return now.Format("20060102") + ".log", nil
	case ROTATION_BY_HOUR:
		return now.Format("2006010210") + ".log", nil
	default:
		return "", ErrNameRotation
	}
}

func GetRotationNameByTime(t time.Time, mode int) (string, error) {
	switch mode {
	case ROTATION_BY_DAY:
		return t.Format("20060102") + ".log", nil
	case ROTATION_BY_HOUR:
		return t.Format("2006010210") + ".log", nil
	default:
		return "", ErrNameRotation
	}
}
