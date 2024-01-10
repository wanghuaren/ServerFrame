package service

import (
	"gameserver/model"
	"gameutils/common"
	"gameutils/common/microuts"
	"gameutils/pbstruct"
	"math/rand"
)

func fightReady(pbDat *pbstruct.CSFightReady, rsp *[]byte, token string, jsonRsp ...*string) {
	_userData := model.CallDBCenter[*pbstruct.MicroUserInfo](common.DB_USERAPI_GetUserDataFromToken, token)
	// _userData := model.GetUserDataFromToken(token)
	_result := pbstruct.SCFightReady{}
	for i := 0; i < int(pbDat.Num); i++ {
		_sc := pbstruct.SCUserInfo{}
		// _sc.Id = rand.Int31n(99)

		randNumMin := _userData.Weapon - common.DefaultWeaponID
		randNumMin = min(randNumMin, 3)
		var randNumMax int32 = 0

		if len(model.StaticTables.Weapon) > 0 {
			randNumMax = model.StaticTables.Weapon[len(model.StaticTables.Weapon)-1].WeaponId - _userData.Weapon
			randNumMax = min(randNumMax, 3)
		}

		_sc.Weapon = _userData.Weapon

		if rand.Intn(10) > 5 && randNumMax > 0 {
			_sc.Weapon = _userData.Weapon + rand.Int31n(randNumMax)
		} else if randNumMin > 0 {
			_sc.Weapon = _userData.Weapon - rand.Int31n(randNumMin)
		}

		// _sc.Money = rand.Int63n(1000)
		// _sc.Strength = rand.Int63n(int64(model.StaticTables.TTili[0].Tilizhi))
		_result.OtherPlay = append(_result.OtherPlay, &_sc)
	}
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, 0, 0, token)
	*rsp = mstBuf
}

func fightResult(pbDat *pbstruct.CSFightResult, rsp *[]byte, token string, jsonRsp ...*string) {
	_result := pbstruct.SCFightResult{}
	_resultErr := 0

	_award := &pbstruct.Reward{}
	_award.GameRk = pbDat.RankNum
	microTransKC := model.SendMicroDB(common.DB_ORDER_FIND, _award, []string{"GameRk"})
	if microTransKC == nil {
		_resultErr = 1
	} else {
		_userData := model.CallDBCenter[*pbstruct.MicroUserInfo](common.DB_USERAPI_GetUserDataFromToken, token)
		// _userData := model.GetUserDataFromToken(token)
		if _userData == nil {
			_resultErr = 2
		} else {
			var retInv = microTransKC.(*pbstruct.RewardResult)
			// microuts.DeMTOriginDatFromMTDat(microTransKC, &retInv)
			if len(retInv.Data) < 1 {
				LogDebug("没有奖励", pbDat.RankNum)
			} else {
				_dat := retInv.Data[0]
				_oldMoney := _userData.Money
				model.CallDBCenter[bool](common.DB_USERAPI_AddUserItem, _userData.Id, int32(_dat.ItemId), int32(_dat.ItemNumber))
				_result.Money = int32(_userData.Money - _oldMoney)
				_result.RankNum = pbDat.RankNum
			}
		}
	}
	model.PrintBackData(pbDat, &_result)
	mstBuf := microuts.GetMTBufFromOriginDat(&_result, pbstruct.MicroTransGate_ID, _resultErr, 0, token)
	*rsp = mstBuf
}
