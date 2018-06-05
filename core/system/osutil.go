/*
@Time : 2018/6/4 20:23
@Author : yinsaki
@File : fileutil
*/

package system

import (
	"runtime"
	"log"
)

const (
	OS_LINUX = iota
	OS_X
	OS_WIN
	OS_OTHERS
)

func GetOsFlag() int{
	switch os := runtime.GOOS;os {
	case "darwin":
		return OS_X
	case "linux":
		return OS_LINUX
	case "window":
		return OS_WIN
	default:
		return OS_OTHERS
	}
}

func GetOsEol() string {
	if GetOsFlag() == OS_WIN {
		return "\r\n"
	}
	return "\n"
}

// 异常转错误
func CatchError() {
	if err := recover(); err != nil {
		log.Println("err", err)
	}
}
