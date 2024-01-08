package main

import (
	"baseutils/baseuts"
	"gamedb/model"
	"gamedb/service"
)

func main() {
	defer baseuts.ChkRecover()

	model.InitModel()
	service.InitService()

	menuOrderStr := "输入命令:\n"
	menuOrderStr += "a:备份MySQL数据库\n"
	menuOrderStr += "b:删除MySQL,把Redis导入MySQL\n"
	menuOrderStr += "c:所有玩家数据存入Redis\n"
	menuOrderStr += "d:显示所有玩家内存中数据\n"
	menuOrderStr += "gc:手动GC\n"
	menuOrderStr += "exit:退出服务\n"
	menuOrderStr += "按Enter键返回菜单\n"
	baseuts.Log(menuOrderStr)
	baseuts.Log("本机IP:" + model.LocalIP)
	baseuts.CommandLine(map[string]func(){
		"a": model.BackUp,
		"b": model.Cache2DB,
		"c": model.FlushAllData2RedisNow,
		"d": model.ShwoAllUserInMemony,
	}, menuOrderStr)
}
