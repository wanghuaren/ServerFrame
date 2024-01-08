package pbuts

import (
	"gameutils/pbstruct"
	"time"
)

func GetCPFromProtoID(protoType int32, buffes ...[]byte) interface{}

func GetCPFromProtoName(_protoName string) interface{}

func GetProtoIDNameFromCP(clientPBDat interface{}) (int32, string)

func GetCTFromCPDat(pbDat interface{}) *pbstruct.ClientTrans {
	var result = pbstruct.ClientTrans{}
	if pbDat != nil {
		result.Protobuff, result.Id = ProtoMarshal(pbDat)
	}
	result.Time = int32(time.Now().Unix())
	return &result
}

func GetCTBufFromCPDat(pbDat interface{}) []byte {
	result := GetCTFromCPDat(pbDat)
	b, _ := ProtoMarshal(result)
	return b
}

func DeCTBuf2CTdat(buf []byte, result ...interface{}) *pbstruct.ClientTrans

func CallProtoEvt(callFunc interface{}, param interface{}, rsp *[]byte, token string, jsonRsp ...*string)

func DeCT2CPDat(clientTrans *pbstruct.ClientTrans) interface{} {
	ret := GetCPFromProtoID(clientTrans.Id, clientTrans.Protobuff)
	return ret
}

func ProtoMarshal(pbDat interface{}) ([]byte, int32)

func ProtoUnMarshal(pbBuf []byte, dstPbDat interface{}) int32
