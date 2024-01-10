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
	ProtoDispatch.AddProtoFunc(pbstruct.CSStaticTab_ID, dbconfig)

	ProtoDispatch.AddProtoFunc(pbstruct.CSUserInfo_ID, userinfo)
	ProtoDispatch.AddProtoFunc(pbstruct.CSBag_ID, userBag)
	ProtoDispatch.AddProtoFunc(pbstruct.CSTask_ID, userTask)
	ProtoDispatch.AddProtoFunc(pbstruct.CSFightReady_ID, fightReady)
	ProtoDispatch.AddProtoFunc(pbstruct.CSFightResult_ID, fightResult)
	ProtoDispatch.AddProtoFunc(pbstruct.CSSkipAD_ID, skipad)

	ProtoDispatch.AddProtoFunc(pbstruct.CSComposeItem_ID, composeItem)
	ProtoDispatch.AddProtoFunc(pbstruct.CSMoveItem_ID, moveItem)
	ProtoDispatch.AddProtoFunc(pbstruct.CSBuyItem_ID, buyItem)
	ProtoDispatch.AddProtoFunc(pbstruct.CSDelItem_ID, delItem)

	ProtoDispatch.AddProtoFunc(pbstruct.CSPay_ID, payment)

	ProtoDispatch.AddProtoFunc(pbstruct.CSSaveError_ID, errlog)
	ProtoDispatch.AddProtoFunc(pbstruct.CSSyncServerTime_ID, syncServerTime)

	ProtoDispatch.AddProtoFunc(pbstruct.CSGMOrder_ID, gamemaster)
}
