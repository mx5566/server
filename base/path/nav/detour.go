package detour

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"unsafe"
)

// 右手坐标系
type NavMeshSetHeader struct {
	magic      int32
	version    int32
	numTiles   int32
	params     DtNavMeshParams
	boundsMinX float32
	boundsMinY float32
	boundsMinZ float32
	boundsMaxX float32
	boundsMaxY float32
	boundsMaxZ float32
}

type NavMeshTileHeader struct {
	tileRef  DtTileRef
	dataSize int32
}

const NAVMESHSET_MAGIC int32 = int32('M')<<24 | int32('S')<<16 | int32('A')<<8 | int32('T')
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

func (d *Detour) Load(path string) error {
	if n, ok := navMeshs[path]; !ok {
		navMesh, err := d.LoadStaticMesh(path)
		if err != nil {
			return err
		}

		d.navMesh = navMesh
	} else {
		d.navMesh = n
	}

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
