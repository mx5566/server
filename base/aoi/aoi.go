package aoi

import "github.com/mx5566/server/base"

type AOI struct {
	x          base.Coord // x方向格子坐标
	y          base.Coord // y方向格子坐标
	unit       ICallBack  // aoi目标对象
	EntityData interface{}
	realNode   *AoiNode // 实际的节点指针
	ID         int64
}

type ICallBack interface {
	OnEnterAoi(aoi *AOI)
	OnLeaveAoi(aoi *AOI)
}

func InitAoi(aoi *AOI, u ICallBack, d interface{}) {
	aoi.unit = u
	aoi.EntityData = d
}
