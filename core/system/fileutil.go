/*
@Time : 2018/6/5 11:23 
@Author : yinsaki
@File : fileutil
*/
package system

import (
	"os"
	"fmt"
	"path/filepath"
)

func GetFileSize(file string) int64 {
	fileInfo, err := os.Stat(file)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	return fileInfo.Size()
}

func IsFileExit(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

func SplitDirFile(path string) (string, string) {
	return filepath.Dir(path), filepath.Base(path)
}