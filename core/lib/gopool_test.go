/*
@Time : 2018/6/7 18:01 
@Author : yinsaki
@File : gopool_test
*/
package lib

import (
	"testing"
	"time"
	"fmt"
	"runtime"
)

func TestNewPool(t *testing.T) {
	p := NewPool(5,10).Start()
	defer p.StopAll()
	for i:= 0; i <100; i ++ {
		count := i
		p.AddJob(func(){
			time.Sleep(time.Second)
			fmt.Println(count)
		})
	}
	time.Sleep(2 * time.Second)
}

func TestPool_EnableWaitForAll(t *testing.T) {
	p := NewPool(3,10).Start().EnableWaitForAll(true)

	for i:= 0; i <2; i ++ {
		count := i
		p.AddJob(func(){
			time.Sleep(time.Second)
			fmt.Println(count)
		})
	}
	p.WaitForAll()
	p.StopAll()
}


func TestPool_StopAll(t *testing.T) {
	rnum := runtime.NumGoroutine()
	p:= NewPool(5, 10).Start().EnableWaitForAll(true)

	defer func() {
		p.StopAll()
		if rnum != runtime.NumGoroutine() {
			t.Error("goroutine not stop", rnum, runtime.NumGoroutine())
		}else {
			t.Log("goroutine stoped")
		}
	}()
	fmt.Println("goroutine num", rnum)

	for i:= 0; i < 10; i++ {
		count := i
		p.AddJob(func(){
			time.Sleep(time.Second)
			fmt.Println(count)
		})
	}

	p.WaitForAll()
}
