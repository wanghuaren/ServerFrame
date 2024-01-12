package model

import (
	"baseutils/baseuts"
	"gameutils/common"
	"gameutils/common/dbuts"
	"gameutils/common/microuts"
	"gameutils/common/pbuts"
	"gameutils/gameuts"
	"gameutils/pbstruct"
	"os"

	_ "gameutils/dbcenter"
	_ "unsafe"
)

var microServer gameuts.IMicroServer
var dbMicroClient gameuts.IMicroClient
var gateMicroClient gameuts.IMicroClient

func initMicro() {
	// discovery := gameuts.Discovery(Conf.String("micro_db_name"), Conf.String("consul_address"))
	// Log(discovery)
	dbMicroClient = gameuts.FactoryClient(Conf.String("micro_db_host"), Conf.String("micro_db_port"), Conf.String("micro_db_name"), Conf.String("micro_name"))
	// discovery := gameuts.Discovery(Conf.String("micro_gate_name"), Conf.String("consul_address"))
	// Log(discovery)
	gateMicroClient = gameuts.FactoryClient(Conf.String("micro_gate_host"), Conf.String("micro_gate_port"), Conf.String("micro_gate_name"), Conf.String("micro_name"))

	microServer = gameuts.FactoryServer(Conf.String("host"), Conf.String("micro_port"), Conf.String("micro_name"))
	microServer.RegisiterFunc(Conf.String("micro_db_name"), receiveMicro)
	microServer.RegisiterFunc(Conf.String("micro_gate_name"), receiveMicro)
	go func() {
		defer baseuts.ChkRecover()
		microServer.Run()
		os.Exit(0)
	}()
}

func SendMicroDB(act string, dbData interface{}, relatedField ...[]string) interface{} {
	if dbData == nil {
		return nil
	}
	microDB := dbuts.GetMTDBFromDBDat(dbData, act, relatedField...)
	if microDB == nil {
		Log("microDBBuf nil")
		return nil
	}
	if dbMicroClient == nil {
		Log("dbMicroClient nil")
		return nil
	}
	microTrans := dbMicroClient.CallFunc(microDB)
	if microTrans == nil {
		return nil
	}
	_, _tabProtoResult := dbuts.CreateStructFromDBDat(dbData)
	microuts.DeMTOriginDatFromMTDat(microTrans, _tabProtoResult)
	return _tabProtoResult
}

func SendMicroGate(key string, buff []byte) *pbstruct.MicroTrans {
	var clientPB pbstruct.ClientTrans
	microGate := pbstruct.MicroTransGate{}
	microGate.ClientTrans = pbuts.DeCTBuf2CTdat(buff, &clientPB)
	// b, _ := pbuts.ProtoMarshal(&microGate)
	// if b == nil {
	// 	return nil
	// }
	microTrans := gateMicroClient.CallFunc(&microGate)
	if microTrans == nil {
		return nil
	}
	return microTrans
}

func CallDBCenter[T *pbstruct.UserTab | *pbstruct.MicroUserInfo | *pbstruct.UserTabResult | *pbstruct.MicroUserToken | *pbstruct.MicroUserTokenFixed | *pbstruct.MicroUserLogin | *pbstruct.MicroUserTokenAndUserInfo | bool | string](key string, args ...interface{}) T {
	var _result = sendMicroDBKey(dbMicroClient, key, args...)
	var result T
	if _result != nil {
		result = _result.(T)
	}
	return result
}

// func CallDBCenterBytes[T *pbstruct.MicroUserInfo](key string, args []byte) T {
// 	var _result = sendMicroDBKeyBytes(dbMicroClient, key, args)
// 	var result T
// 	if _result != nil {
// 		result = _result.(T)
// 	}
// 	return result
// }

//go:linkname sendMicroDBKey gameutils/dbcenter.sendMicroDBKey
func sendMicroDBKey(dbMicroClient gameuts.IMicroClient, key string, args ...interface{}) interface{}

//go:linkname sendMicroDBKeyBytes gameutils/dbcenter.sendMicroDBKeyBytes
func sendMicroDBKeyBytes(dbMicroClient gameuts.IMicroClient, key string, args []byte) interface{}

func receiveMicro(pbDat *pbstruct.MicroTrans, rsp *[]byte) {
	// time.Sleep(time.Second * 5)
	if pbDat.ProtoBufType == pbstruct.MicroTransGate_ID {
		var mt pbstruct.MicroTransGate
		microuts.DeMTTypeDatFromMTDat(pbDat, &mt)
		var _tokenDat = CallDBCenter[*pbstruct.MicroUserToken](common.DB_USERAPI_GetUserToken, mt.ClientTrans.Token)
		if _tokenDat == nil && mt.ClientTrans.Id != pbstruct.CSLogin_ID {
			Log("token permission error", mt.ClientTrans)
			mstBuf := microuts.GetMTBufFromOriginDat(nil, pbstruct.MicroTransGate_ID, common.CLIENT_TRANS_TOKEN_ERROR, 0, mt.ClientTrans.Token)
			*rsp = mstBuf
		} else {
			protoCT := pbuts.GetCPFromProtoID(mt.ClientTrans.Id)
			microuts.DeMTOriginDatFromMTDat(pbDat, &protoCT)
			ProtoDispatch.CallProtoFunc(mt.ClientTrans.Id, protoCT, rsp, mt.ClientTrans.Token)
		}
	} else if pbDat.ProtoBufType == pbstruct.MicroTransDB_ID {
	} else {
		rsp = nil
	}
}
