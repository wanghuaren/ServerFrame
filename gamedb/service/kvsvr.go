package service

import (
	"gamedb/model"
	"gameutils/common"
	"gameutils/common/pbuts"
	"gameutils/pbstruct"
)

func getGameDBFunc(mtdb *pbstruct.MicroTransDB) *pbstruct.MicroTransDB {
	LogDebug("getGameDBFunc", mtdb.FindKey, mtdb.FindKeyArgsString, mtdb.FindKeyArgsInt, mtdb.FindKeyArgsBool)
	var _result = &pbstruct.MicroTransDB{}
	_result.FindKey = mtdb.FindKey
	switch mtdb.FindKey {
	case common.DB_USERAPI_GetUserData:
		var _findDat = model.GetUserData(mtdb.FindKeyArgsInt[0])
		if _findDat != nil {
			for i := range _findDat.Bag {
				if _findDat.Bag[i] == nil {
					_findDat.Bag[i] = &pbstruct.BagItem{Pos: -1, ItemId: 0}
				}
			}
			_result.FindResultBytes, _ = pbuts.ProtoMarshal(_findDat)
		} else {
			_result.FindResultBytes = nil
		}
	case common.DB_USERAPI_GetUserDataFromToken:
		var _findDat = model.GetUserDataFromToken(mtdb.FindKeyArgsString[0])
		_result.FindResultBytes, _ = pbuts.ProtoMarshal(_findDat)
	case common.DB_USERAPI_GetUserTokenFixed:
		var _findDat = model.GetUserTokenFixed(mtdb.FindKeyArgsString[0])
		var _resultDat = pbstruct.MicroUserTokenAndUserInfo{}
		if _findDat != nil {
			_resultDat.UserTokenFixed = _findDat
			var _findDatUserST = model.GetUserDataFromToken(_findDat.Token)
			_resultDat.UserInfo = _findDatUserST
		}
		_result.FindResultBytes, _ = pbuts.ProtoMarshal(&_resultDat)
	// case common.DB_USERAPI_GetUserTokenFixedFromToken:
	// 	var _findDat = model.GetUserTokenFixedFromToken(mtdb.FindKeyArgsString[0])
	// 	_result.FindResultBytes, _ = pbuts.ProtoMarshal(_findDat)
	case common.DB_USERAPI_GetUserToken:
		var _findDat = model.GetUserToken(mtdb.FindKeyArgsString[0])
		_result.FindResultBytes, _ = pbuts.ProtoMarshal(_findDat)
	case common.DB_USERAPI_GetUserTokenFromKey:
		var _findDat = model.GetUserTokenFromKey(mtdb.FindKeyArgsString[0])
		_result.FindResultBytes, _ = pbuts.ProtoMarshal(_findDat)
	// case common.DB_USERAPI_SetUserTokenFromKey:
	// 	model.SetUserTokenFromKey(mtdb.FindKeyArgsString[0])
	case common.DB_USERAPI_CleanUserTokenFromKey:
		model.CleanUserTokenFromKey(mtdb.FindKeyArgsString[0])
	// case common.DB_USERAPI_CleanUserIndex:
	// 	model.CleanUserIndex(mtdb.FindKeyArgsString[0])
	case common.DB_USERAPI_HeartJump:
		model.HeartJump(mtdb.FindKeyArgsString[0])
	case common.DB_USERAPI_AddUserItem:
		model.AddUserItem(mtdb.FindKeyArgsInt[0], mtdb.FindKeyArgsInt[1], mtdb.FindKeyArgsInt[2])
	case common.DB_USERAPI_UseupUserItem:
		_result.FindResultBool = model.UseupUserItem(mtdb.FindKeyArgsInt[0], mtdb.FindKeyArgsInt[1], mtdb.FindKeyArgsInt[2])
	case common.DB_TABLEAPI_InitStaticTables:
		var _findDat = model.GetStaticTables()
		_result.FindResultBytes, _ = pbuts.ProtoMarshal(_findDat)
	case common.DB_USERAPI_FinishTask:
		model.FinishTask(mtdb.FindKeyArgsInt[0], mtdb.FindKeyArgsInt[1])
	case common.DB_USERAPI_Login:
		var _errNum, _token, _userInfo = model.Login(mtdb.FindKeyArgsString[0])
		_findDat := pbstruct.MicroUserLogin{ErrNum: int32(_errNum), Token: _token, UserInfo: _userInfo}
		_result.FindResultBytes, _ = pbuts.ProtoMarshal(&_findDat)
	}
	LogDebug("----------", mtdb.FindKey, mtdb.FindResultString, mtdb.FindResultInt, mtdb.FindResultBool)
	return _result
}
