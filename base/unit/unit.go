package unit

import "github.com/mx5566/server/base/aoi"

type UnitType uint8

type Unit struct {
	id       int64    // id
	unitType UnitType // 对象类型
	aoi      aoi.AOI
}

func (u *Unit) Init() {

	aoi.InitAoi(&u.aoi, u)
}

func (u *Unit) OnEnterAoi(aoi *aoi.AOI) {
	// 给aoi发送消息
}

func (u *Unit) OnLeaveAoi(aoi *aoi.AOI) {

}
