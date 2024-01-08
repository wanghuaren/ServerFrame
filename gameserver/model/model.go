package model

import (
	"baseutils/baseuts"
	"gameutils/common/commuts"
	"gameutils/gameuts"
	"net/http"
	"net/url"
)

var Conf gameuts.IConfig = gameuts.InitConf("serverconf.ini", "server")
var _ = baseuts.InitLog("server")
var ProtoDispatch = commuts.ProtoDispatch{}
var Log = baseuts.Log
var LogDebug = baseuts.LogDebug
var ChkErr = baseuts.ChkErr

var IsDebug = baseuts.Debug(Conf.Bool("debug"))

// var StaticTables *pbstruct.SCStaticTab = nil

var ServerInChina = Conf.Bool("server_in_china")

var uri, _ = url.Parse(Conf.String("local_proxy"))
var httpClient = http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyURL(uri),
	},
}

var LocalIP = baseuts.GetLocalIP()

func InitModel() {
	initConsul()
	if !ServerInChina {
		httpClient = http.Client{}
	}

	go getGooglePayAccessToken()
	initMicro()
	// go initStaticTablesMap()

	gameuts.PPROFCheck(Conf.String("host"), Conf.String("pprof_port"))
}
