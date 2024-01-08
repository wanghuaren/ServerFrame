package model

import (
	"gameutils/common/commuts"
	"gameutils/common/pbuts"
	"reflect"
)

func PrintBackData(pbDat interface{}, backDat interface{}) {
	if IsDebug {
		pbDatRef := reflect.ValueOf(pbDat)
		if pbDatRef.Kind() == reflect.Ptr && pbDatRef.IsNil() {
			Log("接收 nil")
		} else {
			_id, _name := pbuts.GetProtoIDNameFromCP(pbDat)
			Log("接收", _id, _name, commuts.Struct2Map(pbDat))
		}

		backDatRef := reflect.ValueOf(backDat)
		if backDatRef.Kind() == reflect.Ptr && backDatRef.IsNil() {
			Log("返回 nil")
		} else {
			_id, _name := pbuts.GetProtoIDNameFromCP(backDat)
			Log("返回", _id, _name, commuts.Struct2Map(backDat))
		}
		Log("------------------------------")
	}
}
