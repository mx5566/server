package aoiskiplist

import "github.com/mx5566/server/base"

type AoiYListProc struct {
	interDist float32 //兴趣范围
	skipList  *SkipList
}

func NewAoiYlistProc(interDist float32, t ComT) *AoiYListProc {
	return &AoiYListProc{
		skipList:  NewSkipList(t),
		interDist: interDist,
	}
}

func (l *AoiYListProc) insert(aoi *AOIProc, x, y base.Coord) {
	aoi.y = y
	aoi.x = x

	// 已经存在的节点重置一下，直接添加
	var nodeNew = aoi.realYNode
	if nodeNew == nil {
		level := randomLevel()
		nodeNew = makeNode(level, aoi)
	} else {
		nodeNew.reset()
	}

	aoi.realYNode = nodeNew

	l.skipList.insert(nodeNew)
}

// 双向链表的删除
func (l *AoiYListProc) remove(nodeDel *SkipListNode) {
	l.skipList.remove(nodeDel)
}

func (l *AoiYListProc) move(node *SkipListNode, xNew, yNew base.Coord) {
	l.remove(node)
	l.insert(node.aoi, xNew, yNew)
}

func (l *AoiYListProc) getRange(node *SkipListNode) (*SkipListNode, *SkipListNode) {
	minPos, maxPos := node.aoi.x-base.Coord(l.interDist), node.aoi.x+base.Coord(l.interDist)

	return l.skipList.getRange(minPos, maxPos)
}

func (l *AoiYListProc) mark(aoiNode *SkipListNode) {
	minPos, maxPos := aoiNode.aoi.y-base.Coord(l.interDist), aoiNode.aoi.y+base.Coord(l.interDist)
	pre := aoiNode.backward
	for pre != nil && pre.aoi.y >= minPos {
		pre.aoi.mark += 1
		pre = pre.backward
	}

	/*next := aoiNode.level[0].forward
	for next != nil && next.aoi.y < maxPos {
		next.aoi.mark += 1
		next = next.level[0].forward
	}*/

	next := aoiNode.forward
	for next != nil && next.aoi.y < maxPos {
		next.aoi.mark += 1
		next = next.forward
	}
}

func (l *AoiYListProc) clearMark(node *SkipListNode) {
	minPos := node.aoi.y - base.Coord(l.interDist) // 左边界
	maxPos := node.aoi.y + base.Coord(l.interDist) // 又边界

	pre := node.backward
	for pre != nil && pre.aoi.y >= minPos {
		pre.aoi.mark = 0
		pre = pre.backward
	}

	/*	next := node.level[0].forward
		for next != nil && next.aoi.y < maxPos {
			next.aoi.mark = 0
			next = next.level[0].forward
		}*/

	next := node.forward
	for next != nil && next.aoi.y < maxPos {

		next.aoi.mark = 0
		next = next.forward
	}
}
