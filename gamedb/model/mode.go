package model

import (
	"baseutils/baseuts"
	"gameutils/common/commuts"
	"gameutils/gameuts"
)

var Conf gameuts.IConfig = gameuts.InitConf("dbconf.ini", "database")
var _ = baseuts.InitLog("db")

var ProtoDispatch = commuts.ProtoDispatch{}
var Log = baseuts.Log
var LogDebug = baseuts.LogDebug
var ChkErr = baseuts.ChkErr

var IsDebug = baseuts.Debug(Conf.Bool("debug"))

const UserSerialNumberBase = 100000

const UserDropLineTimeout int64 = 60 * 2

// var StaticTables *pbstruct.SCStaticTab = nil

var LocalIP = baseuts.GetLocalIP()

func InitModel() {
	initConsul()
	initDB()
	initRedis()
	dB2Cache()
	initMicro()
	go CheckUserFlushData2Redis()

	gameuts.PPROFCheck(Conf.String("host"), Conf.String("pprof_port"))
}
