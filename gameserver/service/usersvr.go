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
	_dat := model.CallDBCenter[*pbstruct.MicroUserInfo](common.DB_USERAPI_GetUserDataFromToken, token)
	if _dat != nil {
		model.FillUserInfo2Proto(&_result, _dat)
	} else {
		_resultErr = 1
	}
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, 0, token)
	*rsp = mstBuf
}

func userBag(pbDat *pbstruct.CSBag, rsp *[]byte, token string, jsonRsp ...*string) {
	_result := pbstruct.SCBag{}
	_resultErr := 0
	_dat := model.CallDBCenter[*pbstruct.MicroUserInfo](common.DB_USERAPI_GetUserDataFromToken, token)
	if _dat != nil {
		_result.Bag = _dat.Bag[:]
	}
	if _result.Bag == nil {
		_resultErr = 1
	}
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, 0, token)
	*rsp = mstBuf
}

func userTask(pbDat *pbstruct.CSTask, rsp *[]byte, token string, jsonRsp ...*string) {
	_result := pbstruct.SCTask{}
	_resultErr := 0
	_dat := model.CallDBCenter[*pbstruct.MicroUserInfo](common.DB_USERAPI_GetUserDataFromToken, token)
	if _dat != nil {
		_result.Task = _dat.Task

	} else {
		_resultErr = 1
	}
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, 0, token)
	*rsp = mstBuf
}

func skipad(pbDat *pbstruct.CSSkipAD, rsp *[]byte, token string, jsonRsp ...*string) {
	_result := pbstruct.SCSkipAD{}
	_resultErr := 0

	_userTab := &pbstruct.UserTab{}
	_tokenDat := model.CallDBCenter[*pbstruct.MicroUserToken](common.DB_USERAPI_GetUserToken, token)
	if _tokenDat == nil {
		_resultErr = 1
	} else {
		_userTab.Id = _tokenDat.Id
		_userTab.Skipad = pbDat.SkipAD

		microTrans := model.SendMicroDB(common.DB_ORDER_EDIT, _userTab, []string{"Id"}, []string{"Skipad"})
		// var ret pbstruct.UserTabResult
		if microTrans != nil {
			var ret = microTrans.(*pbstruct.UserTabResult)
			// microuts.DeMTOriginDatFromMTDat(microTrans, &ret)
			if len(ret.Data) < 1 {
				_resultErr = 2
			}
		} else {
			_resultErr = 3
		}
	}
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, 0, token)
	*rsp = mstBuf
}

func gamemaster(pbDat *pbstruct.CSGMOrder, rsp *[]byte, token string, jsonRsp ...*string) {
	_result := pbstruct.SCGMOrder{}
	_resultErr := 0
	_userDat := model.CallDBCenter[*pbstruct.MicroUserInfo](common.DB_USERAPI_GetUserDataFromToken, token)
	if _userDat == nil {
		_resultErr = 1
	} else {
		if IsDebug {
			model.GMOrderDO(_userDat, pbDat.OrderNum)
		} else {
			_resultErr = 2
		}
	}
	_result.OrderNum = pbDat.OrderNum
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, 0, token)
	*rsp = mstBuf
}
