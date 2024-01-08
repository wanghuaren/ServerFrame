package main

import (
	"baseutils/baseuts"
	"gamegate/model"
	"gamegate/service"
)

func main() {
	defer baseuts.ChkRecover()

	model.InitMode()
	service.InitService()

	menuOrderStr := "输入命令:\n"
	menuOrderStr += "gc:手动GC\n"
	menuOrderStr += "exit:退出服务\n"
	menuOrderStr += "按Enter键返回菜单\n"
	baseuts.Log(menuOrderStr)
	baseuts.Log("本机IP:" + model.LocalIP)

	baseuts.CommandLine(map[string]func(){}, menuOrderStr)
}
