package service

import (
	"gameserver/model"
	"gameutils/common"
	"gameutils/common/microuts"
	"gameutils/pbstruct"
)

func errlog(pbDat *pbstruct.CSSaveError, rsp *[]byte, token string, jsonRsp ...*string) {
	_result := pbstruct.SCSaveError{}
	_resultErr := 0
	_userInfoData := model.CallDBCenter[*pbstruct.MicroUserInfo](common.DB_USERAPI_GetUserDataFromToken, token)
	// _userInfoData := model.GetUserDataFromToken(token)

	_errTab := pbstruct.ErrorRecord{}

	_errTab.UserId = _userInfoData.Id
	_errTab.ErrContent = pbDat.Content
	_errTab.ErrTime = pbDat.Time
	_errTab.ErrDesc = pbDat.Desc

	microTrans := model.SendMicroDB(common.DB_ORDER_ADD, &_errTab, []string{}, []string{"UserId", "ErrContent", "ErrTime", "ErrDesc"})
	if microTrans == nil {
		_resultErr = 1
	}
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, 0, token)
	*rsp = mstBuf
}
