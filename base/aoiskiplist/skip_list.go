package aoiskiplist

import (
	"github.com/mx5566/server/base"
)

// makeNode 创建一个跳跃表节点
func makeNode(level int16, aoi *AOIProc) *SkipListNode {
	n := &SkipListNode{
		aoi:       aoi,
		level:     make([]*AoiLevel, level),
		neighbors: make(map[*AOIProc]struct{}),
	}

	// 构造层指针
	for i := range n.level {
		n.level[i] = new(AoiLevel)
	}

	return n
}

func randomLevel() int16 {
	level := int16(1)
	for base.Random[float32](0, 1) < 0.25 {
		level++
	}

	if level > MAX_SKIP_LIST_LEVEL {
		return MAX_SKIP_LIST_LEVEL
	}

	return level
}

const MAX_SKIP_LIST_LEVEL int16 = 16

type AoiLevel struct {
	// 节点的子节点
	forward *SkipListNode
}

// aoi十字链表节点
type SkipListNode struct {
	aoi      *AOIProc
	backward *SkipListNode
	forward  *SkipListNode
	level    []*AoiLevel // 节点的层数

	neighbors map[*AOIProc]struct{} // 邻居
	//neighbors map[*SkipListNode]struct{} // 邻居
}

func (s *SkipListNode) reset() {
	s.backward = nil
	s.forward = nil

	for _, v := range s.level {
		v.forward = nil
	}

}

type ComT uint8

// 跳表比较的类型
const (
	Com_X ComT = iota
	Com_y
)

type SkipList struct {
	interDist float32 //兴趣范围
	head      *SkipListNode
	tail      *SkipListNode
	// 节点数 不包括头结点
	length int64
	// 当前的最大有效层级
	level int16

	comT ComT
}

func (t ComT) less(a1, a2 *AOIProc) bool {
	if t == Com_X {
		return a1.x < a2.x
	}

	return a1.y < a2.y
}

func (t ComT) equal(a1, a2 *AOIProc) bool {
	if t == Com_X {
		return a1.x == a2.x
	}

	return a1.y == a2.y
}

func NewSkipList(t ComT) *SkipList {
	return &SkipList{
		level: 1,
		head:  makeNode(MAX_SKIP_LIST_LEVEL, &AOIProc{}),
		comT:  t,
	}
}

//var update = make([]*SkipListNode, MAX_SKIP_LIST_LEVEL)

func (l *SkipList) insert(nodeNew *SkipListNode) {
	update := make([]*SkipListNode, MAX_SKIP_LIST_LEVEL)
	node := l.head

	for i := l.level - 1; i >= 0; i-- {
		for node.level[i].forward != nil &&
			( /*node.level[i].forward.aoi.x < nodeNew.aoi.x*/ l.comT.less(node.level[i].forward.aoi, nodeNew.aoi) ||
				( /*node.level[i].forward.aoi.x == nodeNew.aoi.x*/ l.comT.equal(node.level[i].forward.aoi, nodeNew.aoi) &&
					node.level[i].forward.aoi.ID < nodeNew.aoi.ID)) {
			node = node.level[i].forward
		}

		update[i] = node
	}

	level := int16(len(nodeNew.level))
	if level > l.level {
		for i := l.level; i < level; i++ {
			update[i] = l.head
		}

		l.level = level
	}

	for i := 0; i < int(level); i++ {
		// 链表的插入
		nodeNew.level[i].forward = update[i].level[i].forward

		update[i].level[i].forward = nodeNew

		if i == 0 {
			update[i].forward = nodeNew
		}

	}

	// 设置回退节点
	if update[0] == l.head {
		nodeNew.backward = nil
	} else {
		nodeNew.backward = update[0]
	}

	if nodeNew.level[0].forward != nil {
		nodeNew.level[0].forward.backward = nodeNew
	}

	nodeNew.forward = nodeNew.level[0].forward

	l.length++
}

// 双向链表的删除
func (l *SkipList) remove(nodeDel *SkipListNode) {
	// 储存待删除节点每一层的上一个节点
	update := make([]*SkipListNode, MAX_SKIP_LIST_LEVEL)
	node := l.head
	// 寻找待删除节点
	for i := l.level - 1; i >= 0; i-- {
		for node.level[i].forward != nil &&
			( /*node.level[i].forward.aoi.x < nodeDel.aoi.x*/ l.comT.less(node.level[i].forward.aoi, nodeDel.aoi) ||
				( /*node.level[i].forward.aoi.x == nodeDel.aoi.x*/ l.comT.equal(node.level[i].forward.aoi, nodeDel.aoi) &&
					node.level[i].forward.aoi.ID < nodeDel.aoi.ID)) {
			node = node.level[i].forward
		}
		update[i] = node
	}
	// node在循环中，一直是待删除节点的前一个节点
	// 在最底层的索引处向后移动一位，刚好就是待删除节点
	node = node.level[0].forward

	// 找到该节点
	if node != nil && /*nodeDel.aoi.x == node.aoi.x*/ l.comT.equal(nodeDel.aoi, node.aoi) && nodeDel.aoi.ID == node.aoi.ID {
		l.removeNode(node, update)
		return
	}
	return
}

// 删除找到的节点
func (l *SkipList) removeNode(node *SkipListNode, update []*SkipListNode) {
	// 更新每一层的状态
	for i := int16(0); i < l.level; i++ {
		if update[i].level[i].forward == node {
			update[i].level[i].forward = node.level[i].forward

			if i == 0 {
				update[i].forward = node.level[i].forward
			}
		}
	}
	// 更新后面一个节点的回退指针
	if node.level[0].forward != nil {
		node.level[0].forward.backward = node.backward
	} else {
		l.tail = node.backward
	}
	// 更新跳表中的最大层级
	for l.level > 1 && l.head.level[l.level-1].forward == nil {
		l.level--
	}
	l.length--
}

func (l *SkipList) getRange(minPos, maxPos base.Coord) (*SkipListNode, *SkipListNode) {
	minNode := l.head
	// 找到第一个大于等于minPos的节点
	for level := l.level - 1; level >= 0; level-- {
		for minNode.level[level].forward != nil && minNode.level[level].forward.aoi.x < minPos {
			minNode = minNode.level[level].forward
		}
	}

	minNode = minNode.level[0].forward
	// 大于上边界
	if minNode == nil || minNode.aoi.x >= maxPos {
		return nil, nil
	}

	// 找到最后一个小于maxPos的节点
	maxNode := l.head
	for level := l.level - 1; level >= 0; level-- {
		for maxNode.level[level].forward != nil && maxNode.level[level].forward.aoi.x < maxPos {
			maxNode = maxNode.level[level].forward
		}
	}

	// 小于下边界
	if maxNode == l.head || maxNode.aoi.x < minPos {
		return nil, nil
	}

	return minNode, maxNode
}
