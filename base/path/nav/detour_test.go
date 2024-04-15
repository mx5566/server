package detour

import "testing"

func TestDetour_FindPath(t *testing.T) {
	t.Log("测试开始")

	detour := NewDetour()
	err := detour.Load("./all_tiles_navmesh.bin")
	if err != nil {
		t.Log(err.Error())
		return
	}

	t.Log("文件加载完成")
	_, start := detour.GetRandomPoint()
	_, end := detour.GetRandomPoint()

	t.Logf("启点:%v, 终点:%v\n", start, end)

	pathPoints := [][3]float32{}
	err = detour.FindPath(start, end, &pathPoints)
	if err != nil {
		t.Log(err.Error())
	}

	for _, v := range pathPoints {
		t.Logf("寻路点： %v\n", v)
	}
}
