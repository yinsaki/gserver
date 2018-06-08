/*
@Time : 2018/6/5 18:12 
@Author : yinsaki
@File : main
*/
package main

import "yinsaki/gserver/core/log"

func main() {
	log.SetRollingDaily("./log", "gserver.log")
	log.Debug("1111111111111111 %v", 1223)
	log.Info("1111111111111111 %v", 1223)
	log.Debug("1111111111111111 %v", 1223)
	log.Error("1111111111111111 %v", 1223)
	log.Debug("1111111111111111 %v", 1223)
	log.Debug("1111111111111111 %v", 1223)
}