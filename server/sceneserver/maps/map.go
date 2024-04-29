package maps

import (
	"context"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/aoi"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/rpc3"
	"github.com/mx5566/server/server/sceneserver/player"
	"sync"
)

type MapBaseInfo struct {
	BaseID int64
	Name   string
}

type IMap interface {
}

type Map struct {
	//AOI
	mgr *aoi.AoiManager
	// 寻路

	MapInfo MapBaseInfo

	units sync.Map

	entity.Entity
}

func (m *Map) Init() {
	m.Entity.Init()
	m.Entity.Start()

	// 地图加载

	// 初始化aoi
	m.mgr = aoi.NewAoiAmager(100)

	// 初始化寻路组件

	// 初始化所有的固定怪物/NPC

}

func (m *Map) EnterMap(ctx context.Context, playerID int64) {
	p := player.PLAYERMGR.GetEntity(playerID)
	if p == nil {
		return
	}

	player := p.(*player.Player)

	x := base.Random[float32](10, 500) //randPos(10, 500)
	z := base.Random[float32](10, 500) //randPos(10, 500)

	entity.GEntityMgr.SendMsg(rpc3.RpcHead{ID: playerID}, "Player.SetPosition", x, 0.0, z)
	entity.GEntityMgr.SendMsg(rpc3.RpcHead{ID: playerID}, "Player.SetPosition", x, 0.0, z)

	m.units.Store(&player.Unit, struct {
	}{})

	m.mgr.Enter(&player.Unit.Aoi, player.X, player.Y)
}

func (m *Map) LeaveMap(unit *player.Unit) {
	m.units.Delete(unit)
}
