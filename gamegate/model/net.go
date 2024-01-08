package model

import (
	"gameutils/common/pbuts"
	"gameutils/pbstruct"

	"google.golang.org/protobuf/proto"
)

func SendToServer(b []byte, ipStr string) *pbstruct.ClientTrans {
	clientTrans := &pbstruct.ClientTrans{}
	// var key string
	if IsDebug {
		err := proto.Unmarshal(b, clientTrans)
		if err == nil {
			_requestData := pbuts.GetCPFromProtoID(clientTrans.Id, clientTrans.Protobuff)
			_, _name := pbuts.GetProtoIDNameFromCP(_requestData)

			// key = clientTrans.Key
			LogDebug("请求", ipStr, _name, clientTrans.Token, _requestData)
		} else {
			LogDebug("请求数据错误", ipStr)
		}
	}

	clientTrans = SendMicroServer(b)
	return clientTrans
}
