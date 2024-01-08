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
	ProtoDispatch.AddProtoFunc(pbstruct.CSChangeName_ID, changeName)
	ProtoDispatch.AddProtoFunc(pbstruct.CSUpdateScore_ID, updateScore)
	ProtoDispatch.AddProtoFunc(pbstruct.CSRank_ID, selfRank)

	ProtoDispatch.AddProtoFunc(pbstruct.CSUserInfo_ID, userinfo)
}
