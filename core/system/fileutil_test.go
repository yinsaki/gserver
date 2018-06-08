/*
@Time : 2018/6/6 0:48 
@Author : yinsaki
@File : fileutil_test
*/
package system

import (
	"testing"
	"path/filepath"
	"fmt"
)

func TestSplitDirFile(t *testing.T) {
	fmt.Println(filepath.Dir("./a/b/c"))
	fmt.Println(filepath.Dir("C:/a/b/c"))
	fmt.Println(SplitDirFile("C:/a/b/c"))
}

func TestGetDirFileList(t *testing.T) {
	fmt.Println(GetDirFileList("."))
}