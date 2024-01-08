package service

import (
	"gamegate/model"
	"gameutils/pbstruct"
)

var Conf = model.Conf
var ProtoDispatch = &model.ProtoDispatch
var Log = model.Log
var LogDebug = model.LogDebug
var ChkErr = model.ChkErr

var IsDebug = model.IsDebug

func InitService() {
	initHttp()

	ProtoDispatch.AddProtoFunc(pbstruct.MicroTransGate_ID, receiveMicro)
}
