package aoi

import "github.com/mx5566/server/base"

type AoiXList struct {
	interDist float32 //兴趣范围
	head      *AoiNode
	tail      *AoiNode
}

func NewAoiXlist(interDist float32) *AoiXList {
	return &AoiXList{
		interDist: interDist,
		head:      nil,
		tail:      nil,
	}
}

func (l *AoiXList) insert(node *AoiNode) {
	xPos := node.aoi.x // 获取x轴的坐标点
	if l.head != nil {
		p := l.head // 从头开始遍历 找到第一>=或者到没找到xPos的坐标停止
		for p != nil && p.aoi.x < xPos {
			p = p.xNext
		}
		// now, p == nil or p.coord >= insertCoord
		if p == nil { // if p == nil, insert node at the end of list
			tail := l.tail
			tail.xNext = node
			node.xPrev = tail
			l.tail = node
		} else { // otherwise, p >= node, insert node before p
			prev := p.xPrev
			node.xNext = p
			p.xPrev = node
			node.xPrev = prev

			if prev != nil {
				prev.xNext = node
			} else { // p is the head, so node should be the new head
				l.head = node
			}
		}
	} else { // 空的链表 直接初始化
		l.head = node
		l.tail = node
	}
}

// 双向链表的删除
func (l *AoiXList) remove(node *AoiNode) {
	prev := node.xPrev
	next := node.xNext
	if prev != nil {
		prev.xNext = next
		node.xPrev = nil
	} else {
		l.head = next
	}
	if next != nil {
		next.xPrev = prev
		node.xNext = nil
	} else {
		l.tail = prev
	}
}

func (l *AoiXList) move(node *AoiNode, xOld base.Coord) {
	xNew := node.aoi.x
	if xNew > xOld {
		// moving to next ...
		next := node.xNext
		if next == nil || next.aoi.x >= xNew {
			// 坐标变了，但是在链表中的位置没变
			return
		}
		prev := node.xPrev
		//fmt.Println(1, prev, next, prev == nil || prev.xNext == xzaoi)
		if prev != nil {
			prev.xNext = next // remove xzaoi from list
		} else {
			l.head = next // node is the head, trim it
		}
		next.xPrev = prev

		//fmt.Println(2, prev, next, prev == nil || prev.xNext == next)
		prev, next = next, next.xNext
		for next != nil && next.aoi.x < xNew {
			prev, next = next, next.xNext
			//fmt.Println(2, prev, next, prev == nil || prev.xNext == next)
		}
		//fmt.Println(3, prev, next)
		// no we have prev.X < coord && (next == nil || next.X >= coord), so insert between prev and next
		prev.xNext = node
		node.xPrev = prev
		if next != nil {
			next.xPrev = node
		} else {
			l.tail = node
		}
		node.xNext = next

		//fmt.Println(4)
	} else {
		// moving to prev ...
		prev := node.xPrev
		if prev == nil || prev.aoi.x <= xNew {
			// no need to adjust in list
			return
		}

		next := node.xNext
		if next != nil {
			next.xPrev = prev
		} else {
			l.tail = prev // xzaoi is the head, trim it
		}
		prev.xNext = next // remove xzaoi from list

		next, prev = prev, prev.xPrev
		for prev != nil && prev.aoi.x > xNew {
			next, prev = prev, prev.xPrev
		}
		// no we have next.X > coord && (prev == nil || prev.X <= coord), so insert between prev and next
		next.xPrev = node
		node.xNext = next
		if prev != nil {
			prev.xNext = node
		} else {
			l.head = node
		}
		node.xPrev = prev
	}
}

func (l *AoiXList) mark(node *AoiNode) {
	minPos := node.aoi.x - base.Coord(l.interDist) // 左边界
	maxPos := node.aoi.x + base.Coord(l.interDist) // 又边界

	pre := node.xPrev
	for pre != nil && pre.aoi.x >= minPos {
		pre.markVal += 1
		pre = pre.xPrev
	}

	next := node.xNext
	for next != nil && next.aoi.x < maxPos {
		next.markVal += 1
		next = next.xNext
	}
}

func (l *AoiXList) notifyByMark(node *AoiNode) {
	minPos := node.aoi.x - base.Coord(l.interDist) // 左边界
	maxPos := node.aoi.x + base.Coord(l.interDist) // 又边界

	pre := node.xPrev
	for pre != nil && pre.aoi.x >= minPos {
		if pre.markVal == 2 {
			// 表示需要通知进入
			node.neighbors[pre] = struct{}{}
			// 邻居进入自己的视野
			node.aoi.unit.OnEnterAoi(pre.aoi)
			pre.neighbors[node] = struct{}{}
			// 自己进入邻居的视野
			pre.aoi.unit.OnEnterAoi(node.aoi)
		}

		pre.markVal = 0
		pre = pre.xPrev
	}

	next := node.xNext
	for next != nil && next.aoi.x < maxPos {
		if next.markVal == 2 {
			node.neighbors[next] = struct{}{}
			node.aoi.unit.OnEnterAoi(next.aoi)

			next.neighbors[node] = struct{}{}
			next.aoi.unit.OnEnterAoi(node.aoi)
		}

		next.markVal = 0
		next = next.xNext
	}

}
