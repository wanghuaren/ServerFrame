package microuts

import (
	"baseutils/baseuts"
	"gameutils/common"
	"gameutils/common/dbuts"
	"gameutils/common/pbuts"
	"gameutils/dbstruct"
	"gameutils/pbstruct"

	"reflect"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
)

/*
MICRO_Type_GATE_transType=>1:errNum 2:msgNum

MICRO_Type_DB_transType=>1:act 2:search fields 3:edit fields
*/
func GetMTFromPBDat(pbdat interface{}, transType int32, args ...interface{}) *pbstruct.MicroTrans {
	_postmsg := getMTFromPBDatBase(false, pbdat, transType, args...)
	return _postmsg
}

func GetMTFromPBDatJson(pbdat interface{}, transType int32, args ...interface{}) *pbstruct.MicroTrans {
	_postmsg := getMTFromPBDatBase(true, pbdat, transType, args...)
	return _postmsg
}

func getMTFromPBDatBase(isJson bool, pbdat interface{}, transType int32, args ...interface{}) *pbstruct.MicroTrans {
	_postmsg := pbstruct.MicroTrans{}
	_postmsg.ProtoBufType = transType
	_postmsg.Time = time.Now().Unix()

	switch transType {
	case pbstruct.MicroTransGate_ID:
		var mtg = pbstruct.MicroTransGate{}
		mtg.ClientTrans = pbuts.GetCTFromCPDat(pbdat)

		if len(args) > 0 {
			if errNum, ok := args[0].(int); ok {
				mtg.ClientTrans.Err = int32(errNum)
			} else if errNum, ok := args[0].(int32); ok {
				mtg.ClientTrans.Err = errNum
			} else {
				mtg.ClientTrans.Err = common.CLIENT_TRANS_ERROR_NUM_ERROR
			}
		}

		if len(args) > 1 {
			if msgNum, ok := args[1].(int); ok {
				mtg.ClientTrans.Msg = int32(msgNum)
			} else if msgNum, ok := args[1].(int32); ok {
				mtg.ClientTrans.Msg = msgNum
			} else {
				mtg.ClientTrans.Err = common.CLIENT_TRANS_MSG_NUM_ERROR
			}
		}

		if len(args) > 2 {
			if _token, ok := args[2].(string); ok {
				mtg.ClientTrans.Token = _token
			} else {
				baseuts.Log("ClientTrans Token类型不正确", args[2])
			}
		}
		if isJson {
			str := protojson.Format(&mtg)
			_postmsg.ProtoBufStr = str
		} else {
			_postmsg.ProtoBuf, _ = pbuts.ProtoMarshal(&mtg)
		}
	case pbstruct.MicroTransServer_ID:
	case pbstruct.MicroTransDB_ID:
		if pbdat != nil {
			if isJson {
				if len(args) > 3 {
					_postmsg.ProtoBufStr = dbuts.GetMTDBBufFromDBDatJson(pbdat, args[0].(string), args[1].([]string), args[2].([]string))
				} else if len(args) > 2 {
					_postmsg.ProtoBufStr = dbuts.GetMTDBBufFromDBDatJson(pbdat, args[0].(string), args[1].([]string))
				} else if len(args) > 1 {
					_postmsg.ProtoBufStr = dbuts.GetMTDBBufFromDBDatJson(pbdat, args[0].(string))
				} else {
					_postmsg.ProtoBufStr = dbuts.GetMTDBBufFromDBDatJson(pbdat, "")
				}
			} else {
				if len(args) > 3 {
					_postmsg.ProtoBuf = dbuts.GetMTDBBufFromDBDat(pbdat, args[0].(string), args[1].([]string), args[2].([]string))
				} else if len(args) > 2 {
					_postmsg.ProtoBuf = dbuts.GetMTDBBufFromDBDat(pbdat, args[0].(string), args[1].([]string))
				} else if len(args) > 1 {
					_postmsg.ProtoBuf = dbuts.GetMTDBBufFromDBDat(pbdat, args[0].(string))
				} else {
					_postmsg.ProtoBuf = dbuts.GetMTDBBufFromDBDat(pbdat, "")
				}
			}
		}
	}

	return &_postmsg
}

func GetMTFromOriginDat(originData interface{}, microType int32, args ...interface{}) *pbstruct.MicroTrans {
	_trans := getMTBufFromOriginDatBase(false, originData, microType, args...)
	return _trans.(*pbstruct.MicroTrans)
}

func GetMTBufFromOriginDat(originData interface{}, microType int32, args ...interface{}) []byte {
	_trans := GetMTFromOriginDat(originData, microType, args...)
	b, _ := pbuts.ProtoMarshal(_trans)
	return b
}

func GetMTBufFromOriginDatJson(originData interface{}, microType int32, args ...interface{}) string {
	str := getMTBufFromOriginDatBase(true, originData, microType, args...)
	return str.(string)
}

func getMTBufFromOriginDatBase(isJson bool, originData interface{}, microType int32, args ...interface{}) interface{} {
	if isJson {
		_postmsg := GetMTFromPBDatJson(originData, microType, args...)
		str := protojson.Format(_postmsg)
		return str
	} else {
		_postmsg := GetMTFromPBDat(originData, microType, args...)
		return _postmsg
	}
}

func GetMTBufFromMTTypeDat(mtTypeDat interface{}) []byte {
	_postmsg := getMTFromMTTypeDatBase(false, mtTypeDat)
	transBuf, _ := pbuts.ProtoMarshal(_postmsg)
	return transBuf
}

func GetMTJsonFromMTTypeDat(mtTypeDat interface{}) string {
	_postmsg := GetMTFromMTTypeDat(mtTypeDat)
	transJson := protojson.Format(_postmsg)
	return transJson
}

func GetMTFromMTTypeDat(mtTypeDat interface{}) *pbstruct.MicroTrans {
	_postmsg := getMTFromMTTypeDatBase(false, mtTypeDat)
	return _postmsg
}

func GetMTFromMTTypeDatJson(mtTypeDat interface{}) *pbstruct.MicroTrans {
	_postmsg := getMTFromMTTypeDatBase(true, mtTypeDat)
	return _postmsg
}

func getMTFromMTTypeDatBase(isJson bool, mtTypeDat interface{}) *pbstruct.MicroTrans {
	_postmsg := pbstruct.MicroTrans{}
	_postmsg.Time = time.Now().Unix()

	var transBuf []byte
	var transBufStr string
	var err error
	switch mtTypeDat := mtTypeDat.(type) {
	case *pbstruct.MicroTransDB:
		_postmsg.ProtoBufType = pbstruct.MicroTransDB_ID
		if isJson {

			transBufStr = protojson.Format(mtTypeDat)
		} else {
			transBuf, _ = pbuts.ProtoMarshal(mtTypeDat)
		}
	case *pbstruct.MicroTransGate:
		_postmsg.ProtoBufType = pbstruct.MicroTransGate_ID
		if isJson {
			transBufStr = protojson.Format(mtTypeDat)
		} else {
			transBuf, _ = pbuts.ProtoMarshal(mtTypeDat)
		}
	case *pbstruct.MicroTransServer:
		_postmsg.ProtoBufType = pbstruct.MicroTransServer_ID
		if isJson {
			transBufStr = protojson.Format(mtTypeDat)
		} else {
			transBuf, _ = pbuts.ProtoMarshal(mtTypeDat)
		}
	}
	baseuts.ChkErr(err)
	if isJson {
		_postmsg.ProtoBufStr = transBufStr
	} else {
		_postmsg.ProtoBuf = transBuf
	}
	return &_postmsg
}

func GetMTBufFromMTDat(mtDat *pbstruct.MicroTrans) []byte {
	b, _ := pbuts.ProtoMarshal(mtDat)
	return b
}

func DeMTFromMTBuf(buff []byte) *pbstruct.MicroTrans {
	_postmsg := pbstruct.MicroTrans{}
	pbuts.ProtoUnMarshal(buff, &_postmsg)
	return &_postmsg
}

func DeMTTypeDatFromMTDat[T pbstruct.MicroType](mtDat *pbstruct.MicroTrans, result *T) {
	var _retValue reflect.Value
	switch mtDat.ProtoBufType {
	case pbstruct.MicroTransGate_ID:
		var mtg = pbstruct.MicroTransGate{}
		pbuts.ProtoUnMarshal(mtDat.ProtoBuf, &mtg)
		_retValue = reflect.ValueOf(&mtg).Elem()
	case pbstruct.MicroTransServer_ID:
	case pbstruct.MicroTransDB_ID:
		var mtdb = pbstruct.MicroTransDB{}
		pbuts.ProtoUnMarshal(mtDat.ProtoBuf, &mtdb)
		_retValue = reflect.ValueOf(&mtdb).Elem()
	}
	_value := reflect.ValueOf(result).Elem()
	_value.Set(_retValue)
}

func DeMTOriginDatFromMTDat[T pbstruct.MicroType | pbstruct.ProtoType | dbstruct.TableType | interface{}](mtDat *pbstruct.MicroTrans, result T) {
	var _retValue reflect.Value
	switch mtDat.ProtoBufType {
	case pbstruct.MicroTransGate_ID:
		var mtg = pbstruct.MicroTransGate{}
		pbuts.ProtoUnMarshal(mtDat.ProtoBuf, &mtg)
		_ret := pbuts.GetCPFromProtoID(mtg.ClientTrans.Id, mtg.ClientTrans.Protobuff)
		_retValue = reflect.ValueOf(_ret).Elem()
	case pbstruct.MicroTransServer_ID:
	case pbstruct.MicroTransDB_ID:
		var mtdb = pbstruct.MicroTransDB{}
		pbuts.ProtoUnMarshal(mtDat.ProtoBuf, &mtdb)
		_, _ret := dbuts.CreateTabAndListFromTabID(mtdb.Id, &mtdb)
		_retValue = reflect.ValueOf(_ret).Elem()
	}
	if _retValue.Kind() != reflect.Invalid {
		_value := reflect.ValueOf(result).Elem()
		_value.Set(_retValue)
	}
}
