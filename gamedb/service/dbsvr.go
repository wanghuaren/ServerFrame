package service

import (
	"gamedb/model"
	"gameutils/common"
	"gameutils/common/dbuts"
	"gameutils/common/microuts"
	"gameutils/pbstruct"
)

func receiveMicro(pbDat *pbstruct.MicroTransDB, rsp *[]byte, token string, jsonRsp ...*string) {
	var _result interface{}
	if pbDat.Act == common.DB_FIND_KEY {
		_result = getGameDBFunc(pbDat)
		if rsp != nil {
			var cb = microuts.GetMTBufFromMTTypeDat(_result)
			*rsp = cb
		} else {
			var cbStr = microuts.GetMTJsonFromMTTypeDat(_result)
			*jsonRsp[0] = cbStr
		}
	} else {
		_dbStruct, _ := dbuts.CreateTabAndListFromTabID(pbDat.Id, pbDat)
		if pbDat.Act == common.DB_ORDER_ADD {
			_result = add(_dbStruct, pbDat.SearchField, pbDat.EditField)
		} else if pbDat.Act == common.DB_ORDER_DEL {
			_result = del(_dbStruct, pbDat.SearchField, pbDat.EditField)
		} else if pbDat.Act == common.DB_ORDER_EDIT {
			_result = edit(_dbStruct, pbDat.SearchField, pbDat.EditField)
		} else if pbDat.Act == common.DB_ORDER_FIND {
			_result = find(_dbStruct, pbDat.SearchField, pbDat.EditField)
		} else if pbDat.Act == common.DB_ORDER_FIND_ADD {
			_result = findAdd(_dbStruct, pbDat.SearchField, pbDat.EditField)
		} else if pbDat.Act == common.DB_ORDER_FIND_TABLE {
			_result = findTable(_dbStruct, pbDat.SearchField, pbDat.EditField)
		} else {
			Log("MicroTransDB 类型错误", pbDat)
		}

		if _result != nil && IsDebug {
			dbuts.PrintDBResult(_result)
		}

		if rsp != nil {
			var cb []byte
			if _result != nil {
				cb = microuts.GetMTBufFromOriginDat(_result, pbstruct.MicroTransDB_ID)
			}
			*rsp = cb
		} else {
			cbStr := ""
			if _result != nil {
				cbStr = microuts.GetMTBufFromOriginDatJson(_result, pbstruct.MicroTransDB_ID)
			}
			*jsonRsp[0] = cbStr
		}
	}
}

func add(_dbStruct interface{}, relatedField ...[]string) interface{} {
	if IsDebug {
		_name, _ := dbuts.GetTabNameIDFromDBDat(_dbStruct)
		LogDebug(common.DB_ORDER_ADD, _name, relatedField)
		dbuts.PrintDBResult(_dbStruct)
	}
	findData := model.RDAdd(_dbStruct, relatedField...)
	return findData
}

func del(_dbStruct interface{}, relatedField ...[]string) interface{} {
	if IsDebug {
		_name, _ := dbuts.GetTabNameIDFromDBDat(_dbStruct)
		LogDebug(common.DB_ORDER_DEL, _name)
		dbuts.PrintDBResult(_dbStruct)
	}
	findData := model.RDDel(_dbStruct, relatedField...)
	return findData
}

func edit(_dbStruct interface{}, relatedField ...[]string) interface{} {
	if IsDebug {
		_name, _ := dbuts.GetTabNameIDFromDBDat(_dbStruct)
		LogDebug(common.DB_ORDER_EDIT, _name, relatedField)
		dbuts.PrintDBResult(_dbStruct)
	}
	findData := model.RDEdit(_dbStruct, relatedField...)
	return findData
}

func find(_dbStruct interface{}, relatedField ...[]string) interface{} {
	if IsDebug {
		_name, _ := dbuts.GetTabNameIDFromDBDat(_dbStruct)
		LogDebug(common.DB_ORDER_FIND, _name, relatedField)
		dbuts.PrintDBResult(_dbStruct)
	}
	findData := model.RDFind(_dbStruct, relatedField...)
	return findData
}

func findAdd(_dbStruct interface{}, relatedField ...[]string) interface{} {
	if IsDebug {
		_name, _ := dbuts.GetTabNameIDFromDBDat(_dbStruct)
		LogDebug(common.DB_ORDER_FIND_ADD, _name, relatedField)
		dbuts.PrintDBResult(_dbStruct)
	}
	findData := model.RDFindAutoAdd(_dbStruct, relatedField...)
	return findData
}

func findTable(_dbStruct interface{}, relatedField ...[]string) interface{} {
	if IsDebug {
		_name, _ := dbuts.GetTabNameIDFromDBDat(_dbStruct)
		LogDebug(common.DB_ORDER_FIND_TABLE, _name, relatedField)
		dbuts.PrintDBResult(_dbStruct)
	}

	findData := model.RDFind(_dbStruct)
	return findData
}
