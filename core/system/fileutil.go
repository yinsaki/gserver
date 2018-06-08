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
	"io/ioutil"
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

func GetDirFileList(path string) []string {
	dir_list, e := ioutil.ReadDir(path)
	if e != nil {
		fmt.Println("read dir failed")
	}

	var fileList []string
	for _,v := range dir_list {
		fileList = append(fileList, v.Name())
	}

	return fileList
}

func GetFullPath(path string) string {
	_,err := ioutil.ReadDir(path)
	if err != nil {
		absPath, _ := filepath.Abs(path)
		return absPath
	}else {
		absPath, _ := filepath.Abs(filepath.Dir(path))
		return absPath
	}

}