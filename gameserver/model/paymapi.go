package model

import (
	"encoding/json"
	"gameutils/pbstruct"
	"io"
	"net/http"
)

type OrderInfo struct {
	Kind                 string `json:"kind"`
	PurchaseTimeMillis   string `json:"purchaseTimeMillis"`   // 支付时间, 毫秒
	PurchaseState        int    `json:"purchaseState"`        // 是否付费: 0 已支付, 1 取消
	ConsumptionState     int    `json:"consumptionState"`     // 是否被消费: 0 未消费, 1 已消费
	DeveloperPayload     string `json:"developerPayload"`     // 开发者透传参数
	OrderId              string `json:"orderId"`              // 谷歌订单号
	AcknowledgementState int    `json:"acknowledgementState"` // 支付类型:  0 测试, 1 真实
}

func PayCheck(pbDat *pbstruct.CSPay) *OrderInfo {
	var _order *OrderInfo
	if ServerInChina {
		_order = &OrderInfo{}
		_order.PurchaseState = 0
	} else {
		_order = getOrder(pbDat.ProductId, pbDat.PurchaseToken, pbDat.PackageName, PayAccessToken)
	}
	return _order
}

// 获取订单信息
func getOrder(productId, token, packageName, accessToken string) *OrderInfo {
	req, err := http.NewRequest("GET", Conf.String("check_url_pay")+"/"+packageName+"/purchases/products/"+productId+"/tokens/"+token, nil)
	if ChkErr(err) {
		return nil
	}

	q := req.URL.Query()
	q.Add("access_token", accessToken)
	req.URL.RawQuery = q.Encode()

	response, err := httpClient.Do(req)
	if ChkErr(err) {
		return nil
	}

	respBytes, err := io.ReadAll(response.Body)
	if ChkErr(err) {
		return nil
	}
	LogDebug("google pay order", string(respBytes))
	result := &OrderInfo{}
	err = json.Unmarshal(respBytes, result)
	if ChkErr(err) || result.OrderId == "" {
		return nil
	}
	return result
}
