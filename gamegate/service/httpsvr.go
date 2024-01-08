package service

import (
	"gamegate/model"
	"gameutils/common/pbuts"
	"gameutils/pbstruct"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var regisiter map[string]func(c *gin.Context) = map[string]func(c *gin.Context){}

func initHttp() {
	for i := 0; i < int(pbstruct.ProtoMax); i++ {
		regisiter[strconv.Itoa(i)] = httpsvr
	}
	model.InitHttp(regisiter, websockerSvr)
}

func httpsvr(c *gin.Context) {
	var bt int64
	if IsDebug {
		bt = time.Now().UnixMilli()
	}
	b, _ := c.GetRawData()
	clientTrans := model.SendToServer(b, c.ClientIP())

	sendToHttpClient(c, clientTrans)

	if IsDebug {
		var allT = time.Now().UnixMilli() - bt
		if allT > 1000 {
			Log("pay time 1000", allT)
		} else if allT > 500 {
			Log("pay time 500", allT)
		} else if allT > 100 {
			Log("pay time 100", allT)
		} else if allT > 50 {
			Log("pay time 50", allT)
		}
	}
}

func sendToHttpClient(c *gin.Context, clientTrans *pbstruct.ClientTrans) {
	if clientTrans.Id < 0 {
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"code": clientTrans.Id,
			"msg":  "请重试",
		})
	} else {
		if IsDebug {
			// clientTrans.Key = key
			_data := pbuts.DeCT2CPDat(clientTrans)
			_, _name := pbuts.GetProtoIDNameFromCP(pbuts.GetCPFromProtoID(clientTrans.Id))
			LogDebug("返回", _name, clientTrans.Token, clientTrans.Err, clientTrans.Msg, len(clientTrans.Protobuff), _data)
		}
		c.ProtoBuf(http.StatusOK, clientTrans)
	}
}
