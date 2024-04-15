package detour

import (
	"errors"
	"fmt"
	"github.com/mx5566/server/base"
	"math"
	"os"
	"sync"
	"unsafe"
)

// 右手坐标系
type NavMeshSetHeader struct {
	magic    int32
	version  int32
	numTiles int32
	params   DtNavMeshParams
}

type NavMeshTileHeader struct {
	tileRef  DtTileRef
	dataSize int32
}

var EPSION float32 = 0.000001

const NAVMESHSET_MAGIC int32 = 'M'<<24 | 'S'<<16 | 'E'<<8 | 'T'
const NAVMESHSET_VERSION int32 = 1

const MAX_QUERY_NODE = 128

var (
	navMeshs map[string]*DtNavMesh
	mutex    sync.Mutex
)

type Detour struct {
	navMesh  *DtNavMesh
	navQuery *DtNavMeshQuery
}

func NewDetour() *Detour {
	return &Detour{}
}

func (d *Detour) Load(path string) error {
	mutex.Lock()
	if n, ok := navMeshs[path]; !ok {
		navMesh, err := d.LoadStaticMesh(path)
		if err != nil {
			return err
		}

		d.navMesh = navMesh
	} else {
		d.navMesh = n

		navMeshs[path] = n
	}
	mutex.Unlock()

	//
	navQuery := DtAllocNavMeshQuery()
	if navQuery == nil {
		return errors.New("alloc naemesh query fail")
	}
	status := navQuery.Init(d.navMesh, MAX_QUERY_NODE)
	if DtStatusFailed(status) {
		return errors.New(fmt.Sprintf("navmesh query init error:%d", status))
	}

	d.navQuery = navQuery
	return nil
}

func (d *Detour) LoadStaticMesh(path string) (*DtNavMesh, error) {
	meshData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	header := (*NavMeshSetHeader)(unsafe.Pointer(&(meshData[0])))
	if header.magic != NAVMESHSET_MAGIC {
		return nil, errors.New("navmesh_magic error")
	}

	if header.version != NAVMESHSET_VERSION {
		return nil, errors.New("naemesh_version error")
	}

	//	fmt.Printf("boundsMin: %f, %f, %f\n", header.boundsMinX, header.boundsMinY, header.boundsMinZ)
	//	fmt.Printf("boundsMax: %f, %f, %f\n", header.boundsMaxX, header.boundsMaxY, header.boundsMaxZ)

	navMesh := DtAllocNavMesh()
	state := navMesh.Init(&header.params)
	if DtStatusFailed(state) {
		return nil, errors.New(fmt.Sprintf("navmesh init error: %v, state:%d", header.params, state))
	}

	d1 := int32(unsafe.Sizeof(*header))
	for i := 0; i < int(header.numTiles); i++ {
		tileHeader := (*NavMeshTileHeader)(unsafe.Pointer(&(meshData[d1])))
		if tileHeader.tileRef == 0 || tileHeader.dataSize == 0 {
			break
		}
		d1 += int32(unsafe.Sizeof(*tileHeader))

		data := meshData[d1 : d1+tileHeader.dataSize]
		state = navMesh.AddTile(data, int(tileHeader.dataSize), DT_TILE_FREE_DATA, tileHeader.tileRef, nil)
		if DtStatusFailed(state) {
			return nil, errors.New(fmt.Sprintf("navmesh add tile  state:%d", state))
		}
		d1 += tileHeader.dataSize
	}
	return navMesh, nil
}

// 比如玩家进入地图，随机找一个点加入
func (d *Detour) GetRandomPoint() (bool, []float32) {
	if d.navQuery == nil {
		return false, []float32{}
	}

	var randomRef DtPolyRef
	var randomPt [3]float32
	filter := DtAllocDtQueryFilter()
	status := d.navQuery.FindRandomPoint(filter, func() float32 {
		return base.Random[float32](0.0, 1.0)
	}, &randomRef, randomPt[:])

	if DtStatusFailed(status) {
		return false, []float32{}
	}

	// 返回结果
	return true, randomPt[:]
}

func (d *Detour) GetNearestPoint(center []float32) (bool, []float32) {
	if d.navQuery == nil {
		return false, []float32{}
	}

	var ref DtPolyRef
	var nearestPt [3]float32 = [3]float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}
	filter := DtAllocDtQueryFilter()
	halfExt := [3]float32{1, 100, 1}
	status := d.navQuery.FindNearestPoly(center, halfExt[:], filter, &ref, nearestPt[:])
	if DtStatusFailed(status) {
		return false, []float32{}
	}

	// 状态成功也可能是没找到的
	if ref == 0 {
		return false, []float32{}
	}

	// 返回点坐标
	return true, nearestPt[:]
}

// 判断某个点是否能够到达
func (d *Detour) IsCanReach(pos []float32) bool {
	ret, _ := d.GetNearestPoint(pos)
	if !ret {
		return false
	}
	return true
}

// 寻路
func (d *Detour) FindPath(startPos, endPos []float32, points *[][3]float32) error {
	if d.navQuery == nil {
		return errors.New("not create navquery")
	}

	halfPoint := [3]float32{1, 100, 1}

	var startRef DtPolyRef
	var endRef DtPolyRef
	filter := DtAllocDtQueryFilter()
	status := d.navQuery.FindNearestPoly(startPos, halfPoint[:], filter, &startRef, nil)
	if DtStatusFailed(status) || startRef == 0 {
		return errors.New(fmt.Sprintf("query nearest start point error state:%d, ref:%d", status, startRef))
	}

	status = d.navQuery.FindNearestPoly(endPos, halfPoint[:], filter, &endRef, nil)
	if DtStatusFailed(status) {
		return errors.New(fmt.Sprintf("query nearest end point error state:%d, ref:%d", status, startRef))
	}

	polys := make([]DtPolyRef, MAX_QUERY_NODE)
	var npolys int
	d.navQuery.FindPath(startRef, endRef, startPos, endPos, filter, polys, &npolys, MAX_QUERY_NODE)
	if npolys <= 0 {
		return errors.New("not find path")
	}

	var epos [3]float32
	DtVcopy(epos[:], endPos)
	if polys[npolys-1] != endRef {
		d.navQuery.ClosestPointOnPoly(polys[npolys-1], endPos, epos[:], nil)
	}

	straightPathFlags := make([]DtStraightPathFlags, MAX_QUERY_NODE)
	straightPathPolys := make([]DtPolyRef, MAX_QUERY_NODE)
	ptlist := make([]float32, 3*MAX_QUERY_NODE)
	var ptCount int
	status = d.navQuery.FindStraightPath(startPos, epos[:], polys, npolys, ptlist, straightPathFlags,
		straightPathPolys, &ptCount, MAX_QUERY_NODE, 0)
	if DtStatusFailed(status) {
		return errors.New(fmt.Sprintf("find straight path error state:%d", status))
	}

	for i := 0; i < ptCount; i++ {
		*points = append(*points, [3]float32{ptlist[3*i+0], ptlist[3*i+1], ptlist[3*i+2]})
	}

	return nil
}

// 判断直线行走是否会碰到障碍物
func (d *Detour) Recast(startPos, endPos []float32, points []float32) bool {
	halfExt := [3]float32{1, 100, 1}
	filter := DtAllocDtQueryFilter()
	var neartestRef DtPolyRef
	var nearestPt [3]float32 = [3]float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}

	status := d.navQuery.FindNearestPoly(startPos, halfExt[:], filter, &neartestRef, nearestPt[:])
	if DtStatusFailed(status) || neartestRef == 0 {
		return false
	}

	var nHitParm float32
	var nHitNormal [3]float32
	var polys [MAX_QUERY_NODE]DtPolyRef
	var nPathCount int
	status = d.navQuery.Raycast(neartestRef, startPos, endPos, filter, &nHitParm, nHitNormal[:], polys[:], &nPathCount, MAX_QUERY_NODE)
	if DtStatusFailed(status) || nHitParm == math.MaxFloat32 {
		return false
	}

	// 撞到障碍物
	// 表示启动就在障碍物里面
	if nHitParm <= EPSION {
		DtVcopy(points, startPos)
	} else if nHitParm > 0.0 && nHitParm < 1.0 {
		DtVlerp(points, startPos, endPos, nHitParm)
	}

	return true
}
