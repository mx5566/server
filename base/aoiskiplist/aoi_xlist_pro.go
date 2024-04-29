package aoiskiplist

import (
	"github.com/mx5566/server/base"
)

type AoiXListProc struct {
	interDist float32 //兴趣范围
	skipList  *SkipList
}

func NewAoiXlistProc(interDist float32, t ComT) *AoiXListProc {
	return &AoiXListProc{
		skipList:  NewSkipList(t),
		interDist: interDist,
	}
}

func (l *AoiXListProc) insert(aoi *AOIProc, x, y base.Coord) {
	aoi.y = y
	aoi.x = x

	// 已经存在的节点重置一下，直接添加
	var nodeNew = aoi.realXNode
	if nodeNew == nil {
		level := randomLevel()
		nodeNew = makeNode(level, aoi)
	} else {
		nodeNew.reset()
	}

	aoi.realXNode = nodeNew
	l.skipList.insert(nodeNew)
}

// 双向链表的删除
func (l *AoiXListProc) remove(nodeDel *SkipListNode) {
	l.skipList.remove(nodeDel)
}

func (l *AoiXListProc) move(node *SkipListNode, xNew, yNew base.Coord) {
	l.remove(node)
	// y坐标用不到
	l.insert(node.aoi, xNew, yNew)
}

func (l *AoiXListProc) getRange(node *SkipListNode) (*SkipListNode, *SkipListNode) {
	minPos, maxPos := node.aoi.x-base.Coord(l.interDist), node.aoi.x+base.Coord(l.interDist)

	return l.skipList.getRange(minPos, maxPos)
}

func (l *AoiXListProc) mark(aoiNode *SkipListNode) {
	minPos, maxPos := aoiNode.aoi.x-base.Coord(l.interDist), aoiNode.aoi.x+base.Coord(l.interDist)

	pre := aoiNode.backward
	for pre != nil && pre.aoi.x >= minPos {
		pre.aoi.mark += 1
		pre = pre.backward
	}

	/*next := aoiNode.level[0].forward
	for next != nil && next.aoi.x < maxPos {
		next.aoi.mark += 1
		next = next.level[0].forward
	}*/

	next := aoiNode.forward
	for next != nil && next.aoi.x < maxPos {
		next.aoi.mark += 1
		next = next.forward
	}
}

func (l *AoiXListProc) notifyByMark(node *SkipListNode) {
	minPos := node.aoi.x - base.Coord(l.interDist) // 左边界
	maxPos := node.aoi.x + base.Coord(l.interDist) // 又边界

	pre := node.backward
	for pre != nil && pre.aoi.x >= minPos {
		if pre.aoi.mark == 2 {
			// 表示需要通知进入
			node.neighbors[pre.aoi] = struct{}{}
			// 邻居进入自己的视野
			node.aoi.unit.OnEnterAoi(pre.aoi)
			pre.neighbors[node.aoi] = struct{}{}
			// 自己进入邻居的视野
			pre.aoi.unit.OnEnterAoi(node.aoi)
		}
		/*		if pre.aoi.mark == 2 {
				// 表示需要通知进入
				node.neighbors[pre] = struct{}{}
				// 邻居进入自己的视野
				node.aoi.unit.OnEnterAoi(pre.aoi)
				pre.neighbors[node] = struct{}{}
				// 自己进入邻居的视野
				pre.aoi.unit.OnEnterAoi(node.aoi)
			}*/

		pre.aoi.mark = 0
		pre = pre.backward
	}

	next := node.forward
	for next != nil && next.aoi.x < maxPos {
		if next.aoi.mark == 2 {
			node.neighbors[next.aoi] = struct{}{}
			node.aoi.unit.OnEnterAoi(next.aoi)

			next.neighbors[node.aoi] = struct{}{}
			next.aoi.unit.OnEnterAoi(node.aoi)
		}

		next.aoi.mark = 0
		next = next.forward
	}
}
