package model

import (
	"gameutils/common"
	"gameutils/pbstruct"
	"time"
)

type HCJL struct {
	Dat  *pbstruct.Hcjlb
	Next *pbstruct.Hcjlb
}

var HCJLBMap = map[int32]HCJL{}
var WeaponMap = map[int32]*pbstruct.Weapon{}

func initStaticTablesMap() {
	Log("拉取静态表")
	StaticTables = CallDBCenter[*pbstruct.SCStaticTab](common.DB_TABLEAPI_InitStaticTables)
	if StaticTables != nil {
		Log("静态表拉取成功")
		for i := 0; i < len(StaticTables.Hcjlb)-1; i++ {
			_dat := StaticTables.Hcjlb[i]
			HCJLBMap[_dat.WeaponId] = HCJL{Dat: _dat, Next: StaticTables.Hcjlb[i+1]}
		}

		for i := range StaticTables.Weapon {
			WeaponMap[StaticTables.Weapon[i].WeaponId] = StaticTables.Weapon[i]
		}
	} else {
		time.Sleep(time.Second)
		initStaticTablesMap()
	}
}
