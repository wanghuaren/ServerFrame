package service

import (
	"gameserver/model"
	"gameutils/common"
	"gameutils/common/microuts"
	"gameutils/common/pbuts"
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
	_datBytes, _ := pbuts.ProtoMarshal(_dat)
	_dat = model.CallDBCenter[*pbstruct.UserTab](common.DB_USERAPI_SetUserDataFromToken, token, _datBytes, "Nickname")
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
	var _dat = &pbstruct.UserTab{UserScore: pbDat.Score}
	_datBytes, _ := pbuts.ProtoMarshal(_dat)
	_dat = model.CallDBCenter[*pbstruct.UserTab](common.DB_USERAPI_GetUserDataFromToken, token, _datBytes, "UserScore")
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
	_result := pbstruct.SCRank{UserRank: []*pbstruct.SCUserInfo{}}
	_resultErr := 0
	_dat := model.CallDBCenter[*pbstruct.UserTabResult](common.DB_USERAPI_GetUserDataRank, token)
	if _dat != nil {
	} else {
		_resultErr = 1
	}
	for i := 0; i < len(_dat.Data); i++ {
		_d := &pbstruct.SCUserInfo{}
		_d.Nickname = _dat.Data[i].Nickname
		_d.Score = _dat.Data[i].UserScore
		_result.UserRank = append(_result.UserRank, _d)
	}
	model.PrintBackData(pbDat, _dat)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, 0, token)
	*rsp = mstBuf
}
