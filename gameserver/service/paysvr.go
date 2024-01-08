package service

// func payment(pbDat *pbstruct.CSPay, rsp *[]byte, token string, jsonRsp ...*string) {
// 	_userDat := model.CallDBCenter[*pbstruct.MicroUserInfo](common.DB_USERAPI_GetUserDataFromToken, token)
// 	_result := pbstruct.SCPay{}
// 	_resultErr := 0
// 	_resultMsg := 25
// 	if strings.ToLower(pbDat.ProductId) == "ad_1" {
// 		if _userDat.Chakan > 0 {
// 			_resultErr = 1
// 		} else {
// 			_userDat.Chakan = int32(time.Now().Unix())
// 		}
// 		_resultMsg = 0
// 	}
// 	if _resultErr == 0 {
// 		var _shopItem *pbstruct.SShop = nil
// 		if model.StaticTables == nil {
// 			_resultErr = 2
// 		} else {
// 			for i := range model.StaticTables.SShop {
// 				if strings.EqualFold(model.StaticTables.SShop[i].Mingcheng, pbDat.ProductId) {
// 					_shopItem = model.StaticTables.SShop[i]
// 					break
// 				}
// 			}
// 			if _shopItem == nil {
// 				_resultErr = 3
// 			} else {
// 				if pbDat.ProductId == "limited_offer_2." {
// 					if !_userDat.IsFirst {
// 						LogDebug("首充", _userDat)
// 					}
// 					_userDat.IsFirst = true
// 				}

// 				if _shopItem.Quguanggao > 0 {
// 					_userDat.SkipAD = true
// 				}

// 				if _resultErr == 0 {
// 					_payResult := model.PayCheck(pbDat)
// 					if _payResult == nil {
// 						_resultErr = 4
// 					} else {
// 						if !model.ServerInChina {
// 							_userOrder := pbstruct.UserOrder{}
// 							_userOrder.OrderId = _payResult.OrderId
// 							_userOrder.OrderStatus = int32(_payResult.PurchaseState)
// 							_createTime, err := strconv.ParseInt(_payResult.PurchaseTimeMillis, 10, 64)
// 							if !ChkErr(err) {
// 								_userOrder.OrderTime = int32(_createTime / 1000)
// 							}
// 							_userOrder.Packagename = pbDat.PackageName
// 							_userOrder.ProductId = pbDat.ProductId
// 							_userOrder.ProductQuantity = pbDat.Quantity

// 							_userOrder.UserId = _userDat.Id

// 							microTrans := model.SendMicroDB(common.DB_ORDER_ADD, &_userOrder, nil, []string{"OrderId", "OrderStatus", "OrderTime", "Packagename", "ProductId", "ProductQuantity", "UserId"})
// 							if microTrans == nil {
// 								baseuts.LogF("订单记录失败", &_userOrder)
// 							}
// 						} else {
// 							_resultMsg = 25
// 						}
// 					}
// 				}
// 			}
// 		}
// 		if _resultErr == 0 {
// 			model.CallDBCenter[bool](common.DB_USERAPI_AddUserItem, _userDat.Id, _shopItem.Chaopiao, _shopItem.Count)
// 		}
// 	}

// 	_result.ProductId = pbDat.ProductId

// 	if _resultErr > 0 {
// 		_resultMsg = 26
// 	}
// 	model.PrintBackData(pbDat, &_result)
// 	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, _resultMsg, token)
// 	*rsp = mstBuf
// }
