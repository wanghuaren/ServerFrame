package service

import (
	"gameserver/model"
	"gameutils/common"
	"gameutils/common/microuts"
	"gameutils/pbstruct"
)

func composeItem(pbDat *pbstruct.CSComposeItem, rsp *[]byte, token string, jsonRsp ...*string) {
	_result := pbstruct.SCComposeItem{DelPos: -1}
	_resultErr := 0
	_resultMsg := 0
	_userInfoData := model.CallDBCenter[*pbstruct.MicroUserInfo](common.DB_USERAPI_GetUserDataFromToken, token)
	if _userInfoData == nil {
		_resultErr = 1
	} else {
		_userID := _userInfoData.Id

		dstBagItem := _userInfoData.Bag[pbDat.ItemDst.Pos]
		srcBagItem := _userInfoData.Bag[pbDat.ItemSrc.Pos]

		if dstBagItem == nil {
			_resultErr = 2
		}
		if srcBagItem == nil {
			_resultErr = 3
		}

		if _resultErr == 0 {
			if srcBagItem.ItemId != dstBagItem.ItemId || srcBagItem.ItemId != pbDat.ItemSrc.ItemId || dstBagItem.ItemId != pbDat.ItemDst.ItemId {
				_resultErr = 4
			}
		}
		if _resultErr == 0 {
			if srcBagItem.Pos == dstBagItem.Pos {
				_resultErr = 5
			}
		}
		if _resultErr == 0 {
			var composeItem *pbstruct.Hcjlb = nil
			var composeTarget *pbstruct.Hcjlb = nil

			if _dat, ok := model.HCJLBMap[dstBagItem.ItemId]; ok {
				composeItem = _dat.Dat
				composeTarget = _dat.Next
			}

			if composeItem == nil {
				_resultErr = 6
				_resultMsg = 27
			}

			if _resultErr == 0 {
				_task := model.GetTaskTabDat(composeItem.Renwuid)
				if _task == nil {
					_resultErr = 7
				} else {
					_taskStatus := model.GetUserTask(_userInfoData, composeItem.Renwuid)
					if _taskStatus != nil {
						if !_taskStatus.IsFinish {
							if _task.Times == 2 {
								//不重置
								_userInfoData.Weapon = composeTarget.WeaponId
							}
							model.CallDBCenter[bool](common.DB_USERAPI_AddUserItem, _userID, _task.ItemId, _task.Jianglishuliang)
							model.CallDBCenter[bool](common.DB_USERAPI_FinishTask, _userID, composeItem.Renwuid)
						}

						dstBagItem.ItemId = composeTarget.WeaponId
						_userInfoData.Bag[srcBagItem.Pos] = nil
						_result.ItemDst = dstBagItem
						_result.DelPos = srcBagItem.Pos

						_sc := pbstruct.SCUserInfo{}
						model.FillUserInfo2Proto(&_sc, _userInfoData)
						_result.Userinfo = &_sc

					} else {
						_resultErr = 8
					}
				}
			}
		}
	}
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, _resultMsg, token)
	*rsp = mstBuf
}

func moveItem(pbDat *pbstruct.CSMoveItem, rsp *[]byte, token string, jsonRsp ...*string) {
	_result := pbstruct.SCMoveItem{}
	_resultErr := 0
	_userInfoData := model.CallDBCenter[*pbstruct.MicroUserInfo](common.DB_USERAPI_GetUserDataFromToken, token)

	if _userInfoData == nil {
		_resultErr = 1
	} else {

		bagItem := _userInfoData.Bag[pbDat.ToPos]
		if bagItem == nil {
			bagItem = _userInfoData.Bag[pbDat.ItemDst.Pos]
			if bagItem == nil {
				_resultErr = 2
			} else {
				_userInfoData.Bag[pbDat.ToPos] = bagItem
				_userInfoData.Bag[bagItem.Pos] = nil
				_result.DelPos = bagItem.Pos
				bagItem.Pos = pbDat.ToPos
				_result.ItemDst = bagItem
			}
		} else {
			_resultErr = 3
		}
	}
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, 0, token)
	*rsp = mstBuf
}

func buyItem(pbDat *pbstruct.CSBuyItem, rsp *[]byte, token string, jsonRsp ...*string) {
	_result := pbstruct.SCBuyItem{}
	_resultErr := 0
	_resultMsg := 0
	var _userInfoData = model.CallDBCenter[*pbstruct.MicroUserInfo](common.DB_USERAPI_GetUserDataFromToken, token)
	if _userInfoData == nil {
		_resultErr = 4
	} else {
		//默认买一级道具
		_itemId := common.DefaultWeaponID
		if pbDat.ItemID > 0 {
			_itemId = pbDat.ItemID
		}

		var _weapon *pbstruct.Weapon = model.WeaponMap[_itemId]

		if _weapon != nil && _weapon.Xiaohao > 0 && _weapon.ItemId > 0 {
			ok := model.CallDBCenter[bool](common.DB_USERAPI_UseupUserItem, _userInfoData.Id, _weapon.ItemId, _weapon.Xiaohao)
			if ok {
				_emptyIdex := -1
				for i := range _userInfoData.Bag {
					if _userInfoData.Bag[i] == nil || _userInfoData.Bag[i].ItemId == 0 {
						_emptyIdex = i
						break
					}
				}
				if _emptyIdex < 0 {
					_resultErr = 3
					_resultMsg = 23
				} else {
					_bagItem := &pbstruct.BagItem{}
					_bagItem.ItemId = int32(_itemId)
					_bagItem.Pos = int32(_emptyIdex)

					_result.BagItem = _bagItem
					_result.Money = _userInfoData.Money

					_userInfoData.Bag[_bagItem.Pos] = _bagItem
				}
			} else {
				//固定
				_resultErr = 1
				_resultMsg = 25
			}
		} else {
			_resultErr = 2
		}
	}
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, _resultMsg, token)
	*rsp = mstBuf
}

func delItem(pbDat *pbstruct.CSDelItem, rsp *[]byte, token string, jsonRsp ...*string) {
	_result := pbstruct.SCDelItem{}
	_resultErr := 0
	_resultMsg := 24
	_userInfoData := model.CallDBCenter[*pbstruct.MicroUserInfo](common.DB_USERAPI_GetUserDataFromToken, token)

	if _userInfoData == nil {
		_resultErr = 1
	} else {
		_userInfoData.Bag[pbDat.BagItem.Pos].ItemId = 0
		_result.DelPos = pbDat.BagItem.Pos
	}
	if _resultErr > 0 {
		_resultMsg = 0
	}
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, _resultMsg, token)
	*rsp = mstBuf
}
