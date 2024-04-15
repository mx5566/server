package astar

// Kind* constants refer to tile kinds for input and output.
const (
	// KindPlain (.) is a plain tile with a movement cost of 1.
	GridKindPlain = iota
	GridKindBlock
)

// KindCosts map tile kinds to movement costs.
var GridKindCosts = map[int32]float64{
	GridKindPlain: 1.0,
}

type Grid struct {
	x  int32
	y  int32
	gm *GridMap
}

// PathNeighbors returns the neighbors of the tile, excluding blockers and
// tiles off the edge of the board.
func (t *Grid) PathNeighbors() []Pather {
	neighbors := []Pather{}
	for _, offset := range [][]int32{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	} {
		x := t.x + offset[0]
		y := t.x + offset[1]
		if n := t.gm.Grid(x, y); n != nil {
			if !n.gm.IsHasBlockGrid(x, y) {
				neighbors = append(neighbors, n)
			}
		}
	}
	return neighbors
}

// PathNeighborCost returns the movement cost of the directly neighboring tile.
func (t *Grid) PathNeighborCost(to Pather) float64 {
	return GridKindCosts[to.(*Grid).gm.GetBlockVal(t.x, t.y)]
}

// PathEstimatedCost uses Manhattan distance to estimate orthogonal distance
// between non-adjacent nodes.
func (t *Grid) PathEstimatedCost(to Pather) float64 {
	toT := to.(*Grid)
	absX := toT.x - t.x
	if absX < 0 {
		absX = -absX
	}
	absY := toT.y - t.y
	if absY < 0 {
		absY = -absY
	}
	return float64(absX + absY)
}
