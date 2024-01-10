package model

import (
	"gameutils/pbstruct"
	"time"
)

func GetUserTask(userData *pbstruct.MicroUserInfo, taskID int32) *pbstruct.TaskItem {
	for i := range userData.Task {
		_task := userData.Task[i]
		if _task.TaskId == taskID {
			return _task
		}
	}
	return nil
}

func FillUserInfo2Proto(pbDat *pbstruct.SCUserInfo, structUserDat *pbstruct.MicroUserInfo) {
	pbDat.Weapon = structUserDat.Weapon
	pbDat.Number = structUserDat.Number
	pbDat.Money = structUserDat.Money
	pbDat.Strength = structUserDat.Ticket
	pbDat.SkipAD = structUserDat.SkipAD

	pbDat.Chakan = int32(structUserDat.Chakan)

	structUserDatDate := time.Unix(int64(structUserDat.Chakan), 0)
	if time.Now().YearDay()-structUserDatDate.YearDay() > 0 {
		pbDat.Chakan = 0
		structUserDat.Chakan = 0
	}

	pbDat.FirstPay = structUserDat.IsFirst
	pbDat.RegTime = structUserDat.RegTime
}
