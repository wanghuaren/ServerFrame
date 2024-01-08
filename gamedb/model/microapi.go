package model

import (
	"baseutils/baseuts"
	"gameutils/common/microuts"
	"gameutils/gameuts"
	"gameutils/pbstruct"
	"os"

	"google.golang.org/protobuf/encoding/protojson"
)

var microServer gameuts.IMicroServer
var microClient gameuts.IMicroClient

func initMicro() {
	// discovery := gameuts.Discovery(Conf.String("micro_server_name"), Conf.String("consul_address"))
	// Log(discovery)
	microClient = gameuts.FactoryClient(Conf.String("micro_server_host"), Conf.String("micro_server_port"), Conf.String("micro_server_name"), Conf.String("micro_name"))

	microServer = gameuts.FactoryServer(Conf.String("host"), Conf.String("micro_port"), Conf.String("micro_name"))
	microServer.RegisiterFunc(Conf.String("micro_server_name"), receiveMicro)
	if Conf.String("micro_json_port") != "" {
		microServer.RegisiterFuncJson(Conf.String("micro_server_name")+"json", receiveMicroJson)
		go microServer.StartMicroJson(Conf.String("micro_json_port"))
	}
	go func() {
		defer baseuts.ChkRecover()
		microServer.Run()
		os.Exit(0)
	}()
}

func SendMicroServer(dbData interface{}) *pbstruct.MicroTrans {
	// buff := dbuts.GetMTDBBufFromDBDat(dbData, "")
	// if buff == nil {
	// 	Log("SendMicroServer buff nil")
	// 	return nil
	// }
	microTrans := microClient.CallFunc(dbData)
	if microTrans == nil {
		return nil
	}
	return microTrans
}

func receiveMicro(pbDat *pbstruct.MicroTrans, rsp *[]byte) {
	if pbDat.ProtoBufType == pbstruct.MicroTransDB_ID {
		var mtdb pbstruct.MicroTransDB
		microuts.DeMTTypeDatFromMTDat(pbDat, &mtdb)

		ProtoDispatch.CallProtoFunc(pbstruct.MicroTransDB_ID, &mtdb, rsp, "")
	}
}

func receiveMicroJson(pbDat *pbstruct.MicroTrans, rsp *string) {
	if pbDat.ProtoBufType == pbstruct.MicroTransDB_ID {
		var mtdb = pbstruct.MicroTransDB{}
		err := protojson.Unmarshal([]byte(pbDat.ProtoBufStr), &mtdb)
		if baseuts.ChkErr(err) {
			*rsp = ""
		} else {
			ProtoDispatch.CallProtoFunc(pbstruct.MicroTransDB_ID, &mtdb, nil, "", rsp)
		}

	}
}
