package service

import (
	"gameserver/model"
	"gameutils/common"
	"gameutils/common/microuts"
	"gameutils/pbstruct"
)

func userinfo(pbDat *pbstruct.CSUserInfo, rsp *[]byte, token string, jsonRsp ...*string) {
	_result := pbstruct.SCUserInfo{}
	_resultErr := 0
	_dat := model.CallDBCenter[*pbstruct.UserTab](common.DB_USERAPI_GetUserDataFromToken, token)
	if _dat != nil {
		model.FillUserInfo2Proto(&_result, _dat)
	} else {
		_resultErr = 1
	}
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, 0, token)
	*rsp = mstBuf
}

func changeName(pbDat *pbstruct.CSChangeName, rsp *[]byte, token string, jsonRsp ...*string) {
	_result := pbstruct.SCUserInfo{}
	_resultErr := 0
	var _dat = &pbstruct.UserTab{Nickname: pbDat.Nickname}
	_dat = model.CallDBCenter[*pbstruct.UserTab](common.DB_USERAPI_SetUserDataFromToken, token, _dat, "Nickname")
	if _dat != nil {
		model.FillUserInfo2Proto(&_result, _dat)
	} else {
		_resultErr = 1
	}
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, 0, token)
	*rsp = mstBuf
}

func updateScore(pbDat *pbstruct.CSUpdateScore, rsp *[]byte, token string, jsonRsp ...*string) {
	_result := pbstruct.SCUserInfo{}
	_resultErr := 0
	_dat := model.CallDBCenter[*pbstruct.UserTab](common.DB_USERAPI_GetUserDataFromToken, token)
	if _dat != nil {
		model.FillUserInfo2Proto(&_result, _dat)
	} else {
		_resultErr = 1
	}
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, 0, token)
	*rsp = mstBuf
}

func selfRank(pbDat *pbstruct.CSRank, rsp *[]byte, token string, jsonRsp ...*string) {
	_result := pbstruct.SCUserInfo{}
	_resultErr := 0
	_dat := model.CallDBCenter[*pbstruct.UserTab](common.DB_USERAPI_GetUserDataFromToken, token)
	if _dat != nil {
		model.FillUserInfo2Proto(&_result, _dat)
	} else {
		_resultErr = 1
	}
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, 0, token)
	*rsp = mstBuf
}
