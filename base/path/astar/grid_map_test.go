package astar

import (
	"fmt"
	"github.com/mx5566/server/base"
	"testing"
)

func TestGridMap_Load(t *testing.T) {
	t.Log("A* 寻路")
	gridMap := NewGridMap()
	err := gridMap.Load("./block.json")
	if err != nil {
		return
	}

	t.Log("开始寻路")
	_, vec1 := gridMap.GetRandomPoint()
	_, vec2 := gridMap.GetRandomPoint()

	t.Logf("随机连个点: 起点:%v, 终点: %v\n", vec1, vec2)

	pathList := make([]base.Vector3, 0)
	ret := gridMap.FindPath(vec1, vec2, &pathList)
	if !ret {
		t.Logf("没有找到路")
		return
	}
	// 寻路得到的路点顺序是倒过来的
	for i := 0; i < len(pathList); i++ {
		t.Logf(fmt.Sprintf("寻到的路点:%v\n", pathList[i]))
	}
}
