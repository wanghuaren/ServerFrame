package model

import (
	"baseutils/baseuts"
	"gameutils/pbstruct"
)

// func GetStaticTables() *pbstruct.SCStaticTab {
// 	if StaticTables == nil || IsDebug {
// 		Log("获取静态表")
// 		StaticTables = &pbstruct.SCStaticTab{}

// 		_type := reflect.TypeOf(StaticTables).Elem()
// 		_value := reflect.ValueOf(StaticTables).Elem()

// 		for i := 0; i < _type.NumField(); i++ {
// 			_cpType := _type.Field(i)
// 			if _cpType.Type.Kind() == reflect.Slice && !strings.Contains(_cpType.Name, "unknown") {
// 				_cpValue := _value.Field(i)
// 				_tabName := strings.ToLower(baseuts.CamelCase2UnderScoreCase(_cpType.Name))

// 				dbst, _ := dbuts.CreateTabAndListFromTabName(_tabName)
// 				_tabList := RDFind(dbst)
// 				if _tabList != nil {
// 					_refTabList := reflect.ValueOf(_tabList).Elem()
// 					_refArray := _refTabList.FieldByName("Data")
// 					for l := 0; l < _refArray.Len(); l++ {
// 						value := _refArray.Index(l).Elem()
// 						_pbst := pbuts.GetCPFromProtoName(_tabName)
// 						_dbType := value.Type()
// 						_pbValue := reflect.ValueOf(_pbst).Elem()
// 						for n := 0; n < _dbType.NumField(); n++ {
// 							_fieldName := _dbType.Field(n).Name
// 							if _fieldName == "state" || _fieldName == "sizeCache" || _fieldName == "unknownFields" {
// 								continue
// 							}
// 							_pbSValue := _pbValue.FieldByName(baseuts.UnderScoreCase2CamelCase(_fieldName))

// 							_v := value.Field(n)
// 							if _v.Type().Kind() == reflect.Int {
// 								_vnum := int32(int(_v.Int()))
// 								_v = reflect.ValueOf(_vnum)
// 							}
// 							_pbSValue.Set(_v)
// 						}
// 						_cpValue.Set(reflect.Append(_cpValue, reflect.ValueOf(_pbst)))
// 					}
// 				} else {
// 					Log("静态表 " + _tabName + " 为空")
// 				}
// 			}
// 		}
// 	}
// 	return StaticTables
// }

func CacheOneUserData(userID int32, dbstd ...interface{}) {
	var _userTabList []*pbstruct.UserTab
	if len(dbstd) > 0 {
		_userTabList = dbstd[0].([]*pbstruct.UserTab)
	} else {
		_userTab := &pbstruct.UserTab{}
		_userTab.Id = userID

		_userTabListResult := RDFind(_userTab, []string{"Id"})
		if _userTabListResult != nil {
			_userTabList = _userTabListResult.(*pbstruct.UserTabResult).Data
		}
	}
	if _userTabList == nil {
		baseuts.LogF("查找用户 %v 信息表发生错误", userID)
	} else {
		// microuts.DeMTOriginDatFromMTDat(microTrans, &ret)
		if len(_userTabList) > 0 {
			if len(_userTabList) > 1 {
				Log("CacheOneUserData 唯一用户ID找到多条数据", userID)
			}
			_item := _userTabList[0]
			if _, ok := allUserInfo.Load(int32(_item.Id)); !ok {
				// _userInfo := &pbstruct.MicroUserInfo{}
				// _userInfo := &pbstruct.MicroUserInfo{Bag: []*pbstruct.BagItem{}}
				// for i := 0; i < common.MaxBagCount; i++ {
				// 	_userInfo.Bag = append(_userInfo.Bag, &pbstruct.BagItem{})
				// }
				allUserInfo.Store(int32(_item.Id), _item)
				// allUserInfoFromKey.Store(_item.InventoryId, _userInfo)

				// FillUserInfoFromRedis(_userInfo, _item)

				// if _item.Mission == "" {
				// 	for n := range StaticTables.TTask {
				// 		_taskItem := StaticTables.TTask[n]
				// 		_userInfo.Task = append(_userInfo.Task, &pbstruct.TaskItem{TaskId: _taskItem.Id, IsFinish: false})
				// 	}
				// } else {
				// 	_missionArr := strings.Split(_item.Mission, "#")
				// 	_missionFinishArr := strings.Split(_item.MissionFinish, "#")
				// 	for n := range _missionArr {
				// 		_taskIDStr := _missionArr[n]
				// 		if _taskIDStr != "" {
				// 			_taskID, err := strconv.ParseInt(_taskIDStr, 10, 64)
				// 			if !ChkErr(err) {
				// 				_isFinish := false
				// 				for n := range _missionFinishArr {
				// 					if _missionFinishArr[n] == _taskIDStr {
				// 						_isFinish = true
				// 						break
				// 					}
				// 				}
				// 				_userInfo.Task = append(_userInfo.Task, &pbstruct.TaskItem{TaskId: int32(_taskID), IsFinish: _isFinish})
				// 			}
				// 		}
				// 	}
				// }

				// var _inventoryList []*pbstruct.Inventory
				// if len(dbstd) > 1 {
				// 	_inventoryList = dbstd[1].([]*pbstruct.Inventory)
				// } else {
				// 	_inventoryTab := &pbstruct.Inventory{}
				// 	_inventoryTab.UserId = int32(userID)
				// 	_inventoryListResult := RDFind(_inventoryTab, []string{"User_id"}).(*pbstruct.InventoryResult)
				// 	if _inventoryListResult != nil {
				// 		_inventoryList = _inventoryListResult.Data
				// 	}
				// }
				// if _inventoryList == nil {
				// 	baseuts.LogF("查找用户 %v 库存表发生错误", userID)
				// } else {
				// 	if len(_inventoryList) > 0 {
				// 		for i := range _inventoryList {
				// 			if _inventoryList[i].WeaponId > 0 {
				// 				_bagItem := pbstruct.BagItem{ItemId: int32(_inventoryList[i].WeaponId), Pos: int32(_inventoryList[i].Position)}
				// 				_userInfo.Bag[_bagItem.Pos] = &_bagItem
				// 			}
				// 		}
				// 		Log("初始化用户 %v 库存完成", userID)
				// 	} else {
				// 		baseuts.LogF("找不到用户 %v 库存数据", userID)
				// 	}
				// }
				Log("初始化用户 %v 数据完成", userID)
			} else {
				LogDebug("用户 %v 数据重复初始化", userID)
			}
		} else {
			baseuts.LogF("找不到用户 %v 信息数据", userID)
		}
	}
}
