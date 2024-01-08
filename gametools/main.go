package main

import (
	"baseutils/baseuts"
	"gametools/dbcreatestruct"
	"gametools/execluts"
	"gametools/protobuff"
)

func init() {
	baseuts.InitLog("tools")
}

var cmdOrder = map[string]func(){
	"a": func() {
		execluts.CreateExecl2Proto()
		dbcreatestruct.CreateDBStructStart()
		protobuff.CreatePBStart()
	},
	"b": execluts.ImportExecl2DB,
	"c": execluts.SaveDB2Execl,
}

var menuOrderStr = ""

func main() {
	// cmdOrder["a"]()

	defer baseuts.ChkRecover()

	menuOrderStr += "输入命令:\n"
	menuOrderStr += "a:生成ProroBuff(表结构有更新也要生成一次)\n"
	menuOrderStr += "b:删除MySQL,Execl(.xlsx)导入MySQL\n"
	menuOrderStr += "c:删除Execl(.xlsx),MySQL导出到Execl(.xlsx)\n"
	menuOrderStr += "Enter返回菜单\n"
	baseuts.Log(menuOrderStr)
	baseuts.CommandLine(cmdOrder, menuOrderStr)
}
