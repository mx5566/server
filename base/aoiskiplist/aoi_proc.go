package aoiskiplist

import (
	"github.com/mx5566/server/base"
)

type AOIProc struct {
	x          base.Coord // x方向格子坐标
	y          base.Coord // y方向格子坐标
	mark       int16
	unit       ICallBack // aoi目标对象
	EntityData interface{}
	realXNode  *SkipListNode // 实际的节点指针在X跳表
	realYNode  *SkipListNode // 实际的节点指针在Y跳表
	ID         int64
}

type ICallBack interface {
	OnEnterAoi(aoi *AOIProc)
	OnLeaveAoi(aoi *AOIProc)
}

func InitAoi(aoi *AOIProc, u ICallBack, d interface{}) {
	aoi.unit = u
	aoi.EntityData = d
	//aoi.neighbors = make(map[*SkipListNode]int)
}
