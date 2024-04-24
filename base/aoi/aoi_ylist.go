package aoi

import "github.com/mx5566/server/base"

type AoiYList struct {
	interDist float32 //兴趣范围
	head      *AoiNode
	tail      *AoiNode
}

func NewAoiYlist(interDist float32) *AoiYList {
	return &AoiYList{
		interDist: interDist,
		head:      nil,
		tail:      nil,
	}
}

func (l *AoiYList) insert(node *AoiNode) {
	yPos := node.aoi.y // 获取y轴的坐标点
	if l.head != nil {
		p := l.head // 从头开始遍历 找到第一>=或者到没找到yPos的坐标停止
		for p != nil && p.aoi.y < yPos {
			p = p.yNext
		}
		// now, p == nil or p.coord >= insertCoord
		if p == nil { // if p == nil, insert node at the end of list
			tail := l.tail
			tail.yNext = node
			node.yPrev = tail
			l.tail = node
		} else { // otherwise, p >= node, insert node before p
			prev := p.yPrev
			node.yNext = p
			p.yPrev = node
			node.yPrev = prev

			if prev != nil {
				prev.yNext = node
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
func (l *AoiYList) remove(node *AoiNode) {
	prev := node.yPrev
	next := node.yNext
	if prev != nil {
		prev.yNext = next
		node.yPrev = nil
	} else {
		l.head = next
	}
	if next != nil {
		next.yPrev = prev
		node.yNext = nil
	} else {
		l.tail = prev
	}
}

func (l *AoiYList) move(node *AoiNode, yOld base.Coord) {
	yNew := node.aoi.y
	if yNew > yOld {
		// moving to next ...
		next := node.yNext
		if next == nil || next.aoi.y >= yNew {
			// 坐标变了，但是在链表中的位置没变
			return
		}
		prev := node.yPrev
		//fmt.Println(1, prev, next, prev == nil || prev.yNext == node)
		if prev != nil {
			prev.yNext = next // remove  node from list
		} else {
			l.head = next // node is the head, trim it
		}
		next.yPrev = prev

		//fmt.Println(2, prev, next, prev == nil || prev.yNext == next)
		prev, next = next, next.yNext
		for next != nil && next.aoi.y < yNew {
			prev, next = next, next.yNext
			//fmt.Println(2, prev, next, prev == nil || prev.yNext == next)
		}
		//fmt.Println(3, prev, next)
		// no we have prev.Y < yOld && (next == nil || next.X >= coord), so insert between prev and next
		prev.yNext = node
		node.yPrev = prev
		if next != nil {
			next.yPrev = node
		} else {
			l.tail = node
		}
		node.yNext = next

		//fmt.Println(4)
	} else {
		// moving to prev ...
		prev := node.yPrev
		if prev == nil || prev.aoi.y <= yNew {
			// no need to adjust in list
			return
		}

		next := node.yNext
		if next != nil {
			next.yPrev = prev
		} else {
			l.tail = prev // node is the head, trim it
		}
		prev.yNext = next // remove node from list

		next, prev = prev, prev.yPrev
		for prev != nil && prev.aoi.y > yNew {
			next, prev = prev, prev.yPrev
		}
		// no we have next.y > yOld && (prev == nil || prev.y <= coord), so insert between prev and next
		next.yPrev = node
		node.yNext = next
		if prev != nil {
			prev.yNext = node
		} else {
			l.head = node
		}
		node.yPrev = prev
	}
}

func (l *AoiYList) mark(node *AoiNode) {
	minPos := node.aoi.y - base.Coord(l.interDist) // 左边界
	maxPos := node.aoi.y + base.Coord(l.interDist) // 又边界

	pre := node.yPrev
	for pre != nil && pre.aoi.y >= minPos {
		pre.markVal += 1
		pre = pre.yPrev
	}

	next := node.yNext
	for next != nil && next.aoi.y < maxPos {
		next.markVal += 1
		next = next.yNext
	}
}

func (l *AoiYList) clearMark(node *AoiNode) {
	minPos := node.aoi.y - base.Coord(l.interDist) // 左边界
	maxPos := node.aoi.y + base.Coord(l.interDist) // 又边界

	pre := node.yPrev
	for pre != nil && pre.aoi.y >= minPos {
		pre.markVal = 0
		pre = pre.yPrev
	}

	next := node.yNext
	for next != nil && next.aoi.y < maxPos {
		next.markVal = 0
		next = next.yNext
	}
}
