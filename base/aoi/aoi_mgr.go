package aoi

// aoi十字链表节点
type AoiNode struct {
	aoi *AOI

	neighbors    map[*AoiNode]struct{} // 邻居
	xPrev, xNext *AoiNode              // 当前节点的x轴左边和右边 xPre currNode  xNext
	yPrev, yNext *AoiNode              // 当前节点的y轴下边和上边
	markVal      int                   // 用来标记
}

type AoiManager struct {
	xList *AoiXList
	yList *AoiYList
}

// 在地图场景里面初始化
func NewAoiAmager(interDist int32) *AoiManager {
	return &AoiManager{
		xList: NewAoiXlist(interDist),
		yList: NewAoiYlist(interDist),
	}
}

func (m *AoiManager) Enter(aoi *AOI, x, y int32) {
	aoi.y = y
	aoi.x = x

	node := &AoiNode{
		aoi:       aoi,
		neighbors: make(map[*AoiNode]struct{}),
		xPrev:     nil,
		xNext:     nil,
		yPrev:     nil,
		yNext:     nil,
		markVal:   0,
	}

	aoi.realNode = node

	m.xList.insert(node)
	m.yList.insert(node)

	m.Mark(node)

}

func (m *AoiManager) Leave(aoi *AOI) {
	m.xList.remove(aoi.realNode)
	m.yList.remove(aoi.realNode)

	m.Mark(aoi.realNode)

}

func (m *AoiManager) Move(aoi *AOI, x, y int32) {
	oldX := aoi.x
	oldY := aoi.y
	aoi.x, aoi.y = x, y
	node := aoi.realNode
	if oldX != x {
		m.xList.move(node, oldX)
	}
	if oldY != y {
		m.yList.move(node, oldY)
	}

	m.Mark(node)
}

func (m *AoiManager) Mark(node *AoiNode) {
	m.xList.mark(node)
	m.yList.mark(node)

	for key, _ := range node.neighbors {
		if key.markVal == 2 {
			// 移动之后玩家还在事业范围内
			key.markVal = -1 // 避免下面notify再次通知客户端
		} else {
			// 不在事业范围了，徐亚哦从列表移除掉
			// 从自己的列表移除掉
			delete(node.neighbors, key)
			node.aoi.unit.OnLeaveAoi(key.aoi)
			// 从他的列表移除自己
			delete(key.neighbors, node)
			key.aoi.unit.OnLeaveAoi(node.aoi)
		}
	}

	// 通知所有视野范围的邻居
	m.xList.notifyByMark(node)
	// 清空y轴的标记值
	m.yList.clearMark(node)
}
