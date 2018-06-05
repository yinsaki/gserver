/*
@Time : 2018/6/5 15:12 
@Author : yinsaki
@File : lib
*/
package lib

import (
	"bytes"
	"fmt"
)

func DelimiterConcat(delimiter string, input ... interface{}) string  {
	buffer := bytes.Buffer{}
	l := len(input)
	for i:= 0; i < l; i++ {
		str := fmt.Sprint(input[i])
		buffer.WriteString(str)
		if i < l - 1 {
			buffer.WriteString(delimiter)
		}
	}

	return buffer.String()
}