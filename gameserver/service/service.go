package service

import (
	"gameserver/model"
	"gameutils/pbstruct"
)

var Conf = model.Conf
var ProtoDispatch = &model.ProtoDispatch
var Log = model.Log
var LogDebug = model.LogDebug
var ChkErr = model.ChkErr

var IsDebug = model.IsDebug

func InitService() {
	ProtoDispatch.AddProtoFunc(pbstruct.CSLogin_ID, login)

	ProtoDispatch.AddProtoFunc(pbstruct.CSUserInfo_ID, userinfo)
}
