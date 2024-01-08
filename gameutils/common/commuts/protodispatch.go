package commuts

import (
	"gameutils/common/pbuts"
	"gameutils/pbstruct"
)

type ProtoDispatch struct {
	evtArray [pbstruct.ProtoMax][]interface{}
}

func (e *ProtoDispatch) AddProtoFunc(protoID int32, protoFunc interface{}) {
	var array = e.evtArray[protoID]
	e.evtArray[protoID] = append(array, protoFunc)
}

func (e *ProtoDispatch) RemoveProtoFunc(param interface{}, args ...interface{}) {
	switch param := param.(type) {
	case int32:

		if len(args) > 0 {
			removeFunc(e.evtArray[param], args[0])
		} else {
			e.evtArray[param] = []interface{}{}
		}
	default:
		for i, _ := range e.evtArray {
			removeFunc(e.evtArray[i], param)
		}
	}
}

func removeFunc(funcArray []interface{}, mFunc interface{}) {
	for i, v := range funcArray {
		if v == mFunc {
			funcArray[i] = nil
		}
	}
}

func (e *ProtoDispatch) CallProtoFunc(protoID int32, param interface{}, rsp *[]byte, token string, jsonRsp ...*string) {
	var array = e.evtArray[protoID]
	for _, v := range array {
		pbuts.CallProtoEvt(v, param, rsp, token, jsonRsp...)
	}
}
