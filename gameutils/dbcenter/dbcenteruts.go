package dbcenter

import (
	"baseutils/baseuts"
	"gameutils/common"
	"gameutils/common/pbuts"
	"gameutils/gameuts"
	"gameutils/pbstruct"
)

//lint:ignore U1000 Ignore unused function temporarily for debugging
func sendMicroDBKeyBytes(dbMicroClient gameuts.IMicroClient, key string, args []byte) interface{} {
	return sendMicroDBKeyBase(dbMicroClient, key, nil, nil, nil, args)
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func sendMicroDBKey(dbMicroClient gameuts.IMicroClient, key string, args ...interface{}) interface{} {
	_stringArgs := []string{}
	_intArgs := []int32{}
	_boolArgs := []bool{}
	for i := range args {
		switch _dat := args[i].(type) {
		case int:
			_intArgs = append(_intArgs, int32(_dat))
		case int8:
			_intArgs = append(_intArgs, int32(_dat))
		case int16:
			_intArgs = append(_intArgs, int32(_dat))
		case int32:
			_intArgs = append(_intArgs, _dat)
		case int64:
			_intArgs = append(_intArgs, int32(_dat))
		case string:
			_stringArgs = append(_stringArgs, _dat)
		case bool:
			_boolArgs = append(_boolArgs, _dat)
		}
	}
	return sendMicroDBKeyBase(dbMicroClient, key, _stringArgs, _intArgs, _boolArgs, nil)
}

func sendMicroDBKeyBase(dbMicroClient gameuts.IMicroClient, key string, argsString []string, argsInt []int32, argsBool []bool, argsBytes []byte) interface{} {
	microDB := &pbstruct.MicroTransDB{}
	microDB.Act = common.DB_FIND_KEY
	microDB.FindKey = key
	microDB.FindKeyArgsString = argsString
	microDB.FindKeyArgsInt = argsInt
	microDB.FindKeyArgsBool = argsBool
	microDB.FindKeyArgsBytes = argsBytes

	if dbMicroClient == nil {
		baseuts.Log("dbMicroClient nil")
		return nil
	}
	microTrans := dbMicroClient.CallFunc(microDB)
	if microTrans == nil {
		return nil
	}
	var _mtdb pbstruct.MicroTransDB
	pbuts.ProtoUnMarshal(microTrans.ProtoBuf, &_mtdb)
	switch _mtdb.FindKey {
	case common.DB_USERAPI_GetUserData:
		if _mtdb.FindResultBytes == nil {
			return nil
		}
		var _result = &pbstruct.MicroUserInfo{}
		pbuts.ProtoUnMarshal(_mtdb.FindResultBytes, _result)
		return _result
	case common.DB_USERAPI_GetUserDataFromToken:
		if _mtdb.FindResultBytes == nil {
			return nil
		}
		var _result = &pbstruct.MicroUserInfo{}
		pbuts.ProtoUnMarshal(_mtdb.FindResultBytes, _result)
		return _result
	case common.DB_USERAPI_GetUserTokenFixed:
		if _mtdb.FindResultBytes == nil {
			return nil
		}
		var _result = &pbstruct.MicroUserTokenAndUserInfo{}
		pbuts.ProtoUnMarshal(_mtdb.FindResultBytes, _result)
		return _result
	// case common.DB_USERAPI_GetUserTokenFixedFromToken:
	// 	if _mtdb.FindResultBytes == nil {
	// 		return nil
	// 	}
	// 	var _result = &pbstruct.MicroUserTokenFixed{}
	// 	pbuts.ProtoUnMarshal(_mtdb.FindResultBytes, _result)
	// 	return _result
	case common.DB_USERAPI_GetUserToken:
		if _mtdb.FindResultBytes == nil {
			return nil
		}
		var _result = &pbstruct.MicroUserToken{}
		pbuts.ProtoUnMarshal(_mtdb.FindResultBytes, _result)
		return _result
	case common.DB_USERAPI_GetUserTokenFromKey:
		if _mtdb.FindResultBytes == nil {
			return nil
		}
		var _result = &pbstruct.MicroUserToken{}
		pbuts.ProtoUnMarshal(_mtdb.FindResultBytes, _result)
		return _result
	// case common.DB_USERAPI_SetUserTokenFromKey:

	case common.DB_USERAPI_CleanUserTokenFromKey:

	// case common.DB_USERAPI_CleanUserIndex:

	case common.DB_USERAPI_HeartJump:

	case common.DB_USERAPI_AddUserItem:

	case common.DB_USERAPI_UseupUserItem:
		return _mtdb.FindResultBool

	// case common.DB_TABLEAPI_InitStaticTables:
	// 	if _mtdb.FindResultBytes == nil {
	// 		return nil
	// 	}
	// 	var _result = &pbstruct.SCStaticTab{}
	// 	pbuts.ProtoUnMarshal(_mtdb.FindResultBytes, _result)
	// 	return _result
	case common.DB_USERAPI_FinishTask:

	case common.DB_USERAPI_Login:
		if _mtdb.FindResultBytes == nil {
			return nil
		}
		var _result = &pbstruct.MicroUserLogin{}
		pbuts.ProtoUnMarshal(_mtdb.FindResultBytes, _result)
		return _result
	}
	return nil
}
