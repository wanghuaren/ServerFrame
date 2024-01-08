package service

import (
	"baseutils/baseuts"
	"gamedb/model"
	"gameutils/pbstruct"
)

var Conf = model.Conf
var ProtoDispatch = &model.ProtoDispatch
var Log = model.Log
var LogDebug = baseuts.LogDebug
var ChkErr = model.ChkErr

var IsDebug = model.IsDebug

func InitService() {
	ProtoDispatch.AddProtoFunc(pbstruct.MicroTransDB_ID, receiveMicro)
}
