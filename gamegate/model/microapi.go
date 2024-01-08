package model

import (
	"baseutils/baseuts"
	"gameutils/common"
	"gameutils/common/microuts"
	"gameutils/common/pbuts"
	"gameutils/gameuts"
	"gameutils/pbstruct"
	"os"
)

var microServer gameuts.IMicroServer
var microClient gameuts.IMicroClient

func initMicro() {
	// discovery := gameuts.Discovery(Conf.String("micro_server_name"), Conf.String("consul_address"))
	// Log(discovery)
	microClient = gameuts.FactoryClient(Conf.String("micro_server_host"), Conf.String("micro_server_port"), Conf.String("micro_server_name"), Conf.String("micro_name"))

	microServer = gameuts.FactoryServer(Conf.String("host"), Conf.String("micro_port"), Conf.String("micro_name"))
	microServer.RegisiterFunc(Conf.String("micro_server_name"), receiveMicro)
	go func() {
		defer baseuts.ChkRecover()
		microServer.Run()
		os.Exit(0)
	}()
}

func SendMicroServer(buf []byte) *pbstruct.ClientTrans {
	microGate := pbstruct.MicroTransGate{}
	var clientPB = pbuts.DeCTBuf2CTdat(buf)
	microGate.ClientTrans = clientPB
	microTrans := microClient.CallFunc(&microGate)
	if microTrans == nil {
		return &pbstruct.ClientTrans{Err: common.CLIENT_TRANS_ERROR_MICRO}
	} else {
		var mtg = pbstruct.MicroTransGate{}
		microuts.DeMTTypeDatFromMTDat(microTrans, &mtg)

		if mtg.ClientTrans == nil {
			return &pbstruct.ClientTrans{Err: common.CLIENT_TRANS_ERROR_DECODE}
		}
		return mtg.ClientTrans
	}
}

func receiveMicro(pbDat *pbstruct.MicroTrans, rsp *[]byte) {
	if pbDat.ProtoBufType == pbstruct.MicroTransGate_ID {
		var mtgate pbstruct.MicroTransGate
		microuts.DeMTTypeDatFromMTDat(pbDat, &mtgate)

		ProtoDispatch.CallProtoFunc(pbstruct.MicroTransGate_ID, &mtgate, rsp, "")
	}
}
