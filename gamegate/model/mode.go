package model

import (
	"baseutils/baseuts"
	"gameutils/common/commuts"
	"gameutils/gameuts"
)

var Conf gameuts.IConfig = gameuts.InitConf("gateconf.ini", "gate")
var _ = baseuts.InitLog("gate")
var ProtoDispatch = commuts.ProtoDispatch{}
var Log = baseuts.Log
var LogDebug = baseuts.LogDebug
var ChkErr = baseuts.ChkErr

var IsDebug = baseuts.Debug(Conf.Bool("debug"))

var ServerInChina = Conf.Bool("server_in_china")

var LocalIP = baseuts.GetLocalIP()

func InitMode() {
	initConsul()
	initMicro()
	gameuts.PPROFCheck(Conf.String("host"), Conf.String("pprof_port"))
}
