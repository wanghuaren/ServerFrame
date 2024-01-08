package model

import (
	"baseutils/baseuts"
	"gameutils/common/commuts"
	"gameutils/pbstruct"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// token -> UserToken
var userToken = sync.Map{}

// client_key -> UserToken
var userTokenKey = sync.Map{}

// token -> UserTokenFixed
var userTokenFixedToken = sync.Map{}

// token fixed -> UserTokenFixed
var userTokenFixed = sync.Map{}

// user id -> UserInfoST
var allUserInfo = sync.Map{}

// user key -> UserInfoST
var allUserInfoFromKey = sync.Map{}

func Login(inventoryId string) (int, string, *pbstruct.UserTab) {
	// _result := pbstruct.SCLogin{}
	var _resultErr = 0
	var token = ""

	_userTab := pbstruct.UserTab{}
	_userTab.UserKey = inventoryId
	// _userTab.InventoryId = inventoryId
	_userTab.RegTime = int32(time.Now().Unix())
	_findAddUserTabResult := RDFindAutoAdd(&_userTab, []string{"UserKey"}, []string{"UserKey", "RegTime"})
	if _findAddUserTabResult == nil {
		_resultErr = 2
	} else {
		var retUserTab = _findAddUserTabResult.(*pbstruct.UserTabResult).Data
		// microuts.DeMTOriginDatFromMTDat(microTrans, &ret)
		if len(retUserTab) > 0 {
			_dat := retUserTab[0]

			if len(retUserTab) > 1 {
				Log("login 唯一用户ID找到多条数据", _dat.Id)
			}

			token = CreateUserToken(int64(_dat.Id), _dat.UserKey)

			// _userTabKC := pbstruct.Inventory{}
			// _userTabKC.UserId = _dat.Id
			// _findAddInventoryResult := RDFindAutoAdd(&_userTabKC, []string{"UserId"}, []string{"UserId"})
			// if _findAddInventoryResult == nil {
			// 	_resultErr = 3
			// } else {
			// 	var retInventory = _findAddInventoryResult.(*pbstruct.InventoryResult).Data
			// 	CacheOneUserData(_dat.Id, retUserTab, retInventory)
			_userCacheDat := GetUserData(int32(_dat.Id))
			return _resultErr, token, _userCacheDat
			// }
		} else {
			_resultErr = 4
		}
	}
	return _resultErr, "", nil
}

func CreateUserToken(id int64, inventoryId string) string {
	_ut := pbstruct.MicroUserToken{Id: int32(id), ClientKey: inventoryId}
	fixedStrFromUserInfo := strconv.FormatInt(id, 10) + _ut.ClientKey

	var _tokenFixed = baseuts.GetMd5(fixedStrFromUserInfo)
	var _tokenFixedValue *pbstruct.MicroUserTokenFixed

	if _dat, ok := userTokenFixed.Load(_tokenFixed); ok {
		_tokenFixedValue = _dat.(*pbstruct.MicroUserTokenFixed)
	} else {
		_ut.FixedToken = _tokenFixed
		_tokenFixedValue = &pbstruct.MicroUserTokenFixed{UserIndex: &_ut, Timeout: int32(time.Now().Unix())}
		userTokenFixed.Store(_tokenFixed, _tokenFixedValue)
		userTokenKey.Store(_tokenFixedValue.UserIndex.ClientKey, _tokenFixedValue.UserIndex)
	}

	fixedStrFromUserInfo += strconv.FormatInt(time.Now().UnixMilli(), 10)
	_tokenNew := baseuts.GetMd5(fixedStrFromUserInfo)

	_tokenFixedValue.Token = _tokenNew
	if _tokenFixedValue.TokenOld != "" {
		userToken.Delete(_tokenFixedValue.TokenOld)
		userTokenFixedToken.Delete(_tokenFixedValue.TokenOld)
	}
	userToken.Store(_tokenFixedValue.Token, _tokenFixedValue.UserIndex)
	userTokenFixedToken.Store(_tokenFixedValue.Token, _tokenFixedValue)

	_tokenFixedValue.TokenOld = _tokenFixedValue.Token

	return _tokenFixedValue.Token
}

func GetUserData(userID int32) *pbstruct.UserTab {
	if _, ok := allUserInfo.Load(userID); !ok {
		CacheOneUserData(userID)
	}
	if _dat, ok := allUserInfo.Load(userID); ok {
		if _dat, ok := _dat.(*pbstruct.UserTab); ok {
			return _dat
		}
		return nil
	} else {
		return nil
	}
}

func GetUserDataFromToken(token string) *pbstruct.UserTab {
	if _dat, ok := userToken.Load(token); ok {
		if _dat, ok := _dat.(*pbstruct.MicroUserToken); ok {
			return GetUserData(int32(_dat.Id))
		}
		return nil
	} else {
		return nil
	}
}

func SetUserDataFromToken(token string, pbdat *pbstruct.CSUserInfo, editFields []string) {
	if _dat, ok := userToken.Load(token); ok {
		if _dat, ok := _dat.(*pbstruct.MicroUserToken); ok {
			var userDat = GetUserData(int32(_dat.Id))
			for i := range editFields {
				fieldsName := editFields[i]

				userValue := reflect.ValueOf(userDat).Elem().FieldByName(fieldsName)
				userValueDat := reflect.ValueOf(pbdat).Elem().FieldByName(fieldsName)

				userValue.Set(userValueDat)
			}
		} else {
			baseuts.LogF("SetUserDataFromToken err1", token)
		}
	} else {
		baseuts.LogF("SetUserDataFromToken err2", token)
	}
}

func GetUserTokenFixed(tokenFixed string) *pbstruct.MicroUserTokenFixed {
	if _dat, ok := userTokenFixed.Load(tokenFixed); ok {
		if _dat, ok := _dat.(*pbstruct.MicroUserTokenFixed); ok {
			return _dat
		}
		return nil
	} else {
		return nil
	}
}

func getUserTokenFixedFromToken(token string) *pbstruct.MicroUserTokenFixed {
	if _dat, ok := userTokenFixedToken.Load(token); ok {
		if _dat, ok := _dat.(*pbstruct.MicroUserTokenFixed); ok {
			return _dat
		}
		return nil
	} else {
		return nil
	}
}

func GetUserToken(token string) *pbstruct.MicroUserToken {
	if _dat, ok := userToken.Load(token); ok {
		if _dat, ok := _dat.(*pbstruct.MicroUserToken); ok {
			return _dat
		}
		return nil
	} else {
		return nil
	}
}

func GetUserTokenFromKey(key string) *pbstruct.MicroUserToken {
	if _dat, ok := userTokenKey.Load(key); ok {
		if _dat, ok := _dat.(*pbstruct.MicroUserToken); ok {
			return _dat
		}
		return nil
	} else {
		return nil
	}
}

// func SetUserTokenFromKey(key string) {
// 	userTokenKey.Store(key, &pbstruct.MicroUserToken{Id: -1})
// }

func CleanUserTokenFromKey(key string) {
	userTokenKey.Delete(key)
}

func cleanUserIndex(userTokenFixedStr string) {
	var userIndex = GetUserTokenFixed(userTokenFixedStr)
	if userIndex != nil {
		userToken.Delete(userIndex.Token)
		userTokenFixed.Delete(userTokenFixedStr)
		allUserInfo.Delete(int32(userIndex.UserIndex.Id))
		userTokenKey.Delete(userIndex.UserIndex.ClientKey)
		allUserInfoFromKey.Delete(userIndex.UserIndex.ClientKey)
	} else {
		Log("cleanUserIndex", "清理失败", userTokenFixedStr)
	}
}

func HeartJump(token string) {
	var fixedTokenDat = getUserTokenFixedFromToken(token)
	fixedTokenDat.Timeout = int32(time.Now().Unix())
}

// func AddUserItem(id int32, itemID int32, num int32) {
// 	_userData := GetUserData(id)
// 	switch itemID {
// 	case 20001:
// 		_userData.Money += num
// 	}
// }

// func UseupUserItem(id int32, itemID int32, num int32) bool {
// 	_userData := GetUserData(id)
// 	switch itemID {
// 	case 20001:
// 		if _userData.Money < num {
// 			return false
// 		}
// 		_userData.Money -= num
// 		return true
// 	}
// 	return false
// }

// func FinishTask(userID int32, taskID int32) {
// 	_userData := GetUserData(userID)
// 	for i := range _userData.Task {
// 		_task := _userData.Task[i]
// 		if _task.TaskId == taskID {
// 			_task.IsFinish = true
// 			break
// 		}
// 	}
// }

// func GetTask(userID int32, taskID int32) *pbstruct.TaskItem {
// 	_userData := GetUserData(userID)
// 	for i := range _userData.Task {
// 		_task := _userData.Task[i]
// 		if _task.TaskId == taskID {
// 			return _task
// 		}
// 	}
// 	return nil
// }

func FillUserInfo2Redis(dbST *pbstruct.UserTab, structUserDat *pbstruct.MicroUserInfo) []string {
	dbST.Id = structUserDat.Id
	// dbST.Wuqiid = structUserDat.Weapon
	// dbST.Yonghubianhao = structUserDat.Number
	// dbST.Money = structUserDat.Money
	// dbST.Ticket = structUserDat.Ticket
	// dbST.Skipad = structUserDat.SkipAD

	// dbST.Chakan = structUserDat.Chakan
	// dbST.IsFirst = structUserDat.IsFirst

	return []string{"Wuqiid", "Yonghubianhao", "Money", "Ticket", "Skipad", "Chakan", "IsFirst"}
}

func FillUserInfoFromRedis(structUserDat *pbstruct.MicroUserInfo, dbST *pbstruct.UserTab) {
	structUserDat.Id = int32(dbST.Id)
	structUserDat.Number = int32(UserSerialNumberBase + dbST.Id)
	// structUserDat.ClientKey = dbST.InventoryId

	// structUserDat.Weapon = int32(dbST.Wuqiid)
	// structUserDat.Money = int32(dbST.Money)
	// structUserDat.Ticket = int32(dbST.Ticket)

	// structUserDat.SkipAD = dbST.Skipad
	// structUserDat.Chakan = int32(dbST.Chakan)
	// structUserDat.IsFirst = dbST.IsFirst
}

func CheckUserFlushData2Redis() {
	time.Sleep(time.Second * time.Duration(UserDropLineTimeout))
	userTokenFixed.Range(func(key, value any) bool {
		v := value.(*pbstruct.MicroUserTokenFixed)
		_currTome := time.Now().Unix()
		if _currTome-int64(v.Timeout) > UserDropLineTimeout {
			var userST = GetUserDataFromToken(v.Token)
			Log("用户 " + strconv.Itoa(int(userST.Id)) + " 数据写入Redis")
			flushOneUser2Redis(userST)
			// cleanUserIndex(v.UserIndex.FixedToken)
			Log("用户 " + strconv.Itoa(int(userST.Id)) + " 数据写入Redis完成")
			v.Timeout = int32(_currTome) + 3600
		}
		return true
	})
	CheckUserFlushData2Redis()
}

func ShwoAllUserInMemony() {
	Log("=========显示所有玩家内存中数据Begin==============")
	var count1 = 0
	allUserInfo.Range(func(key, value any) bool {
		Log(key, commuts.Struct2Map(value.(*pbstruct.MicroUserInfo)))
		count1++
		return true
	})
	var count2 = 0
	userToken.Range(func(key, value any) bool {
		Log(key, commuts.Struct2Map(value))
		count2++
		return true
	})
	var count3 = 0
	userTokenFixed.Range(func(key, value any) bool {
		Log(key, commuts.Struct2Map(value))
		count3++
		return true
	})
	var count4 = 0
	userTokenKey.Range(func(key, value any) bool {
		Log(key, commuts.Struct2Map(value))
		count4++
		return true
	})
	var count5 = 0
	allUserInfoFromKey.Range(func(key, value any) bool {
		Log(key, commuts.Struct2Map(value))
		count5++
		return true
	})

	Log("=========显示所有玩家内存中数据End================", count1, count2, count3, count4, count5)
}
func FlushAllData2RedisNow() {
	Log("全部用户数据写入Redis")
	allUserInfo.Range(func(key, value any) bool {
		k := key.(int32)
		v := value.(*pbstruct.UserTab)
		if k != v.Id {
			// continue
		} else {
			flushOneUser2Redis(v)
		}
		return true
	})
	Log("全部用户数据写入Redis完成")
}

func flushOneUser2Redis(userST *pbstruct.UserTab) {
	// _userTab := &pbstruct.UserTab{}
	// _editFields := FillUserInfo2Redis(_userTab, userST)

	_editFields := []string{"Nickname", "UserSex", "UserScore"}
	// for i := range userST.Task {
	// 	_task := userST.Task[i]
	// 	if _task.IsFinish {
	// 		_userTab.MissionFinish += strconv.FormatInt(int64(_task.TaskId), 10) + "#"
	// 	}
	// 	_userTab.Mission += strconv.FormatInt(int64(_task.TaskId), 10) + "#"
	// }

	// if _userTab.MissionFinish != "" {
	// 	_userTab.MissionFinish = _userTab.MissionFinish[:len(_userTab.MissionFinish)-1]
	// }

	// if _userTab.Mission != "" {
	// 	_userTab.Mission = _userTab.Mission[:len(_userTab.Mission)-1]
	// }
	// _editFields = append(_editFields, "MissionFinish", "Mission")
	_editUserTabBack := RDEdit(userST, []string{"Id"}, _editFields)
	if _editUserTabBack == nil {
		Log("FlushAllData A", commuts.Struct2Map(userST))
	} else {
		// for i := range userST.Bag {
		// 	_bagItem := userST.Bag[i]
		// 	if _bagItem == nil || _bagItem.ItemId == 0 {
		// 		continue
		// 	}
		// 	_inventoryTab := &pbstruct.Inventory{}
		// 	_inventoryTab.UserId = userST.Id
		// 	_inventoryTab.Position = _bagItem.Pos
		// 	_inventoryTab.WeaponId = _bagItem.ItemId
		// 	_inventoryTabDat := RDEdit(_inventoryTab, []string{"UserId", "Position"}, []string{"WeaponId"})
		// 	if _inventoryTabDat == nil {
		// 		_inventoryTabAddBack := RDAdd(_inventoryTab, nil, []string{"UserId", "Position", "WeaponId"})
		// 		if _inventoryTabAddBack == nil {
		// 			Log("FlushAllData B", commuts.Struct2Map(_inventoryTab))
		// 		}
		// 	}
		// }

		// _inventoryTab := &pbstruct.Inventory{UserId: userST.Id}
		// _inventoryTabListResult := RDFind(_inventoryTab, []string{"UserId"})
		// if _inventoryTabListResult == nil {
		// 	Log("FlushAllData C", commuts.Struct2Map(_inventoryTab))
		// } else {
		// 	_inventoryTabList := _inventoryTabListResult.(*pbstruct.InventoryResult).Data

		// 	for i := range _inventoryTabList {
		// 		_ret := _inventoryTabList[i]
		// 		_has := false
		// 		for n := range userST.Bag {
		// 			_bagItem := userST.Bag[n]
		// 			if _bagItem == nil || _bagItem.ItemId == 0 {
		// 				continue
		// 			}
		// 			if _ret.Position == _bagItem.Pos {
		// 				_has = true
		// 				break
		// 			}
		// 		}
		// 		if !_has {
		// 			_inventoryTabDelBack := RDDel(_ret, []string{"Id"})
		// 			if _inventoryTabDelBack == nil {
		// 				Log("FlushAllData D", commuts.Struct2Map(_ret))
		// 			}
		// 		}
		// 	}
		// }
	}
}
