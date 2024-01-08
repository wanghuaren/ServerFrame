package service

import (
	"gameserver/model"
	"gameutils/common"
	"gameutils/common/commuts"
	"gameutils/common/microuts"
	"gameutils/pbstruct"
	"strings"
	"sync"
)

var logingNow = sync.Map{}

func login(pbDat *pbstruct.CSLogin, rsp *[]byte, token string, jsonRsp ...*string) {
	_result := pbstruct.SCLogin{}
	var _resultErr int32 = 0

	var userIndex = model.CallDBCenter[*pbstruct.MicroUserToken](common.DB_USERAPI_GetUserTokenFromKey, pbDat.InventoryId)
	if userIndex == nil {
		if _, ok := logingNow.Load(pbDat.InventoryId); ok {
			_resultErr = 5
		} else {
			logingNow.Store(pbDat.InventoryId, true)
			var _checkDat *model.GoogleLoginResult
			if !model.ServerInChina && !strings.Contains(pbDat.Code, "xhhy_") {
				_checkDat = model.GetGoogleLoginDatFromCode(pbDat.Code)
			}
			if _checkDat == nil && !model.ServerInChina && !strings.Contains(pbDat.Code, "xhhy_") {
				_resultErr = 1
			} else {
				findDat := model.CallDBCenter[*pbstruct.MicroUserLogin](common.DB_USERAPI_Login, pbDat.InventoryId)
				if findDat == nil {
					_resultErr = 6
				} else {
					_resultErr = findDat.ErrNum

					token = findDat.Token
					_result.Token = token

					_userST := findDat.UserInfo
					_userST.Nickname = _userST.UserKey
					_sc := pbstruct.SCUserInfo{}
					model.FillUserInfo2Proto(&_sc, _userST)
					_result.Userinfo = &_sc

					// _result.UserBag = &pbstruct.SCBag{Bag: _userST.Bag[:]}
					// _result.StaticTab = model.StaticTables
				}
			}
			logingNow.Delete(pbDat.InventoryId)
		}
	} else if userIndex.Id >= 0 {
		userKeyInfoData := model.CallDBCenter[*pbstruct.MicroUserTokenAndUserInfo](common.DB_USERAPI_GetUserTokenFixed, userIndex.FixedToken)
		if userKeyInfoData != nil && userKeyInfoData.UserTokenFixed != nil {
			userTokenInfo := userKeyInfoData.UserTokenFixed
			token = userTokenInfo.Token
			_userCacheDat := userKeyInfoData.UserInfo
			_result.Token = token
			_sc := pbstruct.SCUserInfo{}
			model.FillUserInfo2Proto(&_sc, _userCacheDat)
			_result.Userinfo = &_sc
			// _result.UserBag = &pbstruct.SCBag{Bag: _userCacheDat.Bag[:]}
			// _result.StaticTab = model.StaticTables
			LogDebug("复用数据", commuts.Struct2Map(&_result))
		} else {
			LogDebug("复用数据出错")
		}
	}
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, 0, token)
	*rsp = mstBuf
}
