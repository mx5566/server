package aoiskiplist

import (
	"github.com/mx5566/server/base"
)

const Epsion base.Coord = 1

type AoiManagerProc struct {
	xList *AoiXListProc
	yList *AoiYListProc
}

// 在地图场景里面初始化
func NewAoiAmager(interDist float32) *AoiManagerProc {
	return &AoiManagerProc{
		xList: NewAoiXlistProc(interDist, Com_X),
		yList: NewAoiYlistProc(interDist, Com_y),
	}
}

func (m *AoiManagerProc) Enter(aoi *AOIProc, x, y base.Coord) {
	m.xList.insert(aoi, x, y)
	m.yList.insert(aoi, x, y)

	m.mark(aoi)
}

func (m *AoiManagerProc) Leave(aoi *AOIProc) {
	m.xList.remove(aoi.realXNode)
	m.yList.remove(aoi.realYNode)

	for k, _ := range aoi.realXNode.neighbors {
		k.unit.OnLeaveAoi(aoi)
		aoi.unit.OnLeaveAoi(k)

		delete(aoi.realXNode.neighbors, k)
		delete(k.realXNode.neighbors, aoi)
	}
}

func (m *AoiManagerProc) Move(aoi *AOIProc, x, y base.Coord) {
	oldX := aoi.x
	oldY := aoi.y
	//aoi.x, aoi.y = x, y

	if oldX != x {
		m.xList.move(aoi.realXNode, x, y)
	}
	if oldY != y {
		aoi.x, aoi.y = oldX, oldY
		m.yList.move(aoi.realYNode, x, y)
	}

	m.mark(aoi)
}

func (m *AoiManagerProc) mark(aoi *AOIProc) {
	m.xList.mark(aoi.realXNode)
	m.yList.mark(aoi.realYNode)

	neignbors := aoi.realXNode.neighbors
	for key, _ := range neignbors {
		if key.mark == 2 {
			// 移动之后玩家还在事业范围内
			key.mark = -1 // 避免下面notify再次通知客户端
		} else {
			// 不在视野范围了，从列表移除掉
			// 从自己的列表移除掉
			delete(neignbors, key)
			aoi.unit.OnLeaveAoi(key)
			// 从他的列表移除自己
			delete(key.realXNode.neighbors, aoi)
			key.unit.OnLeaveAoi(aoi)
		}
	}

	/*
	      	neignbors := aoi.realXNode.neighbors

	   for key, _ := range neignbors {
	   		if key.aoi.mark == 2 {
	   			// 移动之后玩家还在事业范围内
	   			key.aoi.mark = -1 // 避免下面notify再次通知客户端
	   		} else {
	   			// 不在视野范围了，从列表移除掉
	   			// 从自己的列表移除掉
	   			delete(neignbors, key)
	   			aoi.unit.OnLeaveAoi(key.aoi)
	   			// 从他的列表移除自己
	   			delete(key.neighbors, aoi.realXNode)
	   			key.aoi.unit.OnLeaveAoi(aoi)
	   		}
	   	}*/

	/*	for key, _ := range aoi.neighbors {
		if key.aoi.mark == 2 {
			// 移动之后玩家还在事业范围内
			key.aoi.mark = -1 // 避免下面notify再次通知客户端
		} else {
			// 不在视野范围了，从列表移除掉
			// 从自己的列表移除掉
			delete(aoi.neighbors, key)
			aoi.unit.OnLeaveAoi(key.aoi)
			// 从他的列表移除自己
			delete(key.aoi.neighbors, aoi.realXNode)
			key.aoi.unit.OnLeaveAoi(aoi)
		}
	}*/

	// 通知所有视野范围的邻居
	m.xList.notifyByMark(aoi.realXNode)
	// 清空y轴的标记值
	m.yList.clearMark(aoi.realYNode)

}
