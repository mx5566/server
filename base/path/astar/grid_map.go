package astar

import (
	"encoding/json"
	"github.com/mx5566/server/base"
	"os"
)

type GridMapHeader struct {
	Id         int32 `json:"id"`          //地图id
	Width      int32 `json:"width"`       // 地图的宽 多少个单个格子 GridWidth * Width 真正的长度
	Height     int32 `json:"height"`      //地图的高
	GridWidth  int32 `json:"tile_width"`  // 单个格子的单位长度
	GridHeight int32 `json:"tile_height"` // 单个格子的单位高度
}

type GridMapInfo struct {
	Head  GridMapHeader `json:"head"`
	Block [][]int32     `json:"blocks"`
}

func (i *GridMapInfo) GetGridWidth() int32 {
	return i.Head.GridWidth
}

func (i *GridMapInfo) GetGridHeight() int32 {
	return i.Head.GridHeight
}

func (i *GridMapInfo) GetArea() int32 {
	return i.Head.Height * i.Head.Width
}

func (i *GridMapInfo) GetWidth() int32 {
	return i.Head.Width
}

func (i *GridMapInfo) Getheight() int32 {
	return i.Head.Height
}

func (i *GridMapInfo) IsBlock(x, y int32) bool {
	return i.Block[y][x] == GridKindBlock
}

// 格子地图
type GridMap struct {
	grids       []*Grid
	gridMapInfo GridMapInfo
}

func NewGridMap() *GridMap {
	return &GridMap{
		grids: nil,
		gridMapInfo: GridMapInfo{
			Block: make([][]int32, 0),
		},
	}
}

func (g *GridMap) GetBlockVal(x, y int32) int32 {
	return g.gridMapInfo.Block[x][y]
}

func (g *GridMap) Load(path string) error {
	da, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(da, &g.gridMapInfo)
	if err != nil {
		return err
	}

	g.grids = make([]*Grid, g.gridMapInfo.GetArea())
	var height, width = g.gridMapInfo.Getheight(), g.gridMapInfo.GetWidth()
	var i, j int32 = 0, 0
	for i = 0; i < height; i++ {
		for j = 0; j < width; j++ {
			g.grids[i*height+j] = &Grid{
				x:  j,
				y:  i,
				gm: g,
			}
		}
	}

	return nil
}

func (g *GridMap) Grid(x, y int32) *Grid {
	if x < 0 || y < 0 || x >= g.gridMapInfo.GetWidth() || y >= g.gridMapInfo.Getheight() {
		return nil
	}

	if g.gridMapInfo.GetArea() < x*y {
		return nil
	}

	return g.grids[x+y*g.gridMapInfo.GetWidth()]

}

func (g *GridMap) IsHasBlockGrid(x, y int32) bool {
	if x < 0 || y < 0 || x >= g.gridMapInfo.GetWidth() || y >= g.gridMapInfo.Getheight() {
		return false
	}

	if g.gridMapInfo.GetArea() < x*y {
		return false
	}

	return g.gridMapInfo.IsBlock(x, y)
}

func (g *GridMap) FindPath(start, end base.Vector3, points *[]base.Vector3) bool {
	x1, y1 := g.PosToGrid(start)
	x2, y2 := g.PosToGrid(end)

	grid1 := g.Grid(x1, y1)
	grid2 := g.Grid(x2, y2)
	if grid1 == nil || grid2 == nil {
		return false
	}

	pathList, _, found := Path(grid1, grid2)
	if !found {
		return false
	}

	width := g.gridMapInfo.GetGridWidth()
	height := g.gridMapInfo.GetGridHeight()
	for _, v := range pathList {
		grid := v.(*Grid)

		x := float32(grid.x*width) + float32(width)/2
		y := float32(grid.y*height) + float32(height)/2
		*points = append(*points, base.Vector3{X: base.Coord(x), Y: base.Coord(0), Z: base.Coord(y)})
	}

	return true
}

func (g *GridMap) PosToGrid(pos base.Vector3) (x, y int32) {
	x, y = int32(pos.X/base.Coord(g.gridMapInfo.GetGridWidth())), int32(pos.Z/base.Coord(g.gridMapInfo.GetGridHeight()))
	return
}

func (g *GridMap) GetRandomPoint() (bool, base.Vector3) {
	x1, y1 := base.RandomInt32(0, g.gridMapInfo.GetWidth()), base.RandomInt32(0, g.gridMapInfo.Getheight())
	var time = 0
	var x, y int32
	find := true
	for x, y = int32(x1), int32(y1); g.IsHasBlockGrid(x, y); {
		time++
		if time > 100 {
			find = false
			return false, base.Vector3{}
		}

		x, y = base.RandomInt32(0, g.gridMapInfo.GetWidth()), base.RandomInt32(0, g.gridMapInfo.Getheight())
	}

	if !find {
		return false, base.Vector3{}
	}
	return true, base.Vector3{X: base.Coord(x + g.gridMapInfo.GetGridWidth()/2.0), Y: base.Coord(0), Z: base.Coord(y + g.gridMapInfo.GetGridHeight()/2.0)}
}
