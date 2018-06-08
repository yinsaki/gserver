/*
@Time : 2018/6/6 10:54 
@Author : yinsaki
@File : log_test
*/
package log

import (
	"testing"
	"time"
	"math/rand"
)

func TestTrace(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	line := 0
	SetRollingDaily("../../log", "111.log")
	time := time.Tick(time.Second)
	for {
		select {
			case <- time:
				line ++
				Trace(LogLevel(rand.Intn(5)), "the %v log", line)
		}
	}
}