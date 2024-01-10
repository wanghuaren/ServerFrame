package service

import (
	"gameserver/model"
	"gameutils/common/microuts"
	"gameutils/pbstruct"
)

var clientConfBufCache []byte

func dbconfig(pbDat *pbstruct.CSStaticTab, rsp *[]byte, token string, jsonRsp ...*string) {
	if clientConfBufCache == nil || IsDebug {
		_staticTables := model.StaticTables
		model.PrintBackData(pbDat, _staticTables)
		clientConfBufCache = microuts.GetMTBufFromOriginDat(_staticTables, pbstruct.MicroTransGate_ID, 0, 0, token)
	}
	*rsp = clientConfBufCache
}
