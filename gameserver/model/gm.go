package model

import (
	"gameutils/common"
	"gameutils/pbstruct"
)

func GMOrderDO(userDat *pbstruct.MicroUserInfo, orderType int32) {
	switch orderType {
	case 1:
		addMoney(userDat)
	case 2:
		delAllMoney(userDat)
	case 3:
		delAllBag(userDat)
	case 4:
		resetUseWeapon(userDat)
	case 5:
		addWeapon(userDat)
	}
}
func addMoney(userDat *pbstruct.MicroUserInfo) {
	userDat.Money += 1000
}

func delAllMoney(userDat *pbstruct.MicroUserInfo) {
	userDat.Money = 0
}

func delAllBag(userDat *pbstruct.MicroUserInfo) {
	for i := range userDat.Bag {
		userDat.Bag[i] = nil
	}
}

func resetUseWeapon(userDat *pbstruct.MicroUserInfo) {
	userDat.Weapon = common.DefaultWeaponID
	for i := range userDat.Task {
		userDat.Task[i].IsFinish = false
	}
}

func addWeapon(userDat *pbstruct.MicroUserInfo) {
	for i := range userDat.Bag {
		if userDat.Bag[i] == nil {
			userDat.Bag[i] = &pbstruct.BagItem{ItemId: common.DefaultWeaponID + 5, Pos: int32(i)}
			break
		}
	}
}
