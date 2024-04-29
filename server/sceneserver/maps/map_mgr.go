package maps

import (
	"context"
	"github.com/mx5566/server/base/cluster"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/rpc3"
	"reflect"
)

type MapMgr struct {
	entity.Entity
	entity.EntityPool
}

func (m *MapMgr) Init() {
	m.Entity.Init()
	m.EntityPool.InitPool(reflect.TypeOf(Map{}))
	entity.GEntityMgr.RegisterEntity(m)
	m.Entity.Start()

	m.LoadMap()
}

func (m *MapMgr) LoadMap() {
	var i int64
	for i = 10; i < 20; i++ {
		m1 := new(Map)
		m1.SetID(i)
		m1.Init()
		m.AddEntity(m)
	}
}

func (m *MapMgr) EnterMap(ctx context.Context, mapID uint32, playerID int64, gateClusterID, gameClusterID uint32) {
	m1 := m.GetEntity(int64(mapID))
	if m1 == nil {
		//地图不存在
		cluster.GCluster.SendMsg(&rpc3.RpcHead{
			ID:             playerID,
			SrcServerID:    cluster.GCluster.Id(),
			DestServerID:   gameClusterID,
			DestServerType: rpc3.ServiceType_GameServer,
		}, "gameserver<-Player.LoginMapResult", false, mapID)
		return
	}

	entity.GEntityMgr.SendMsg(rpc3.RpcHead{}, "PlayerMgr.LoginMap", mapID, playerID, gateClusterID, gameClusterID)
}

func (m *MapMgr) LeaveMap(ctx context.Context, mapID uint32, playerID int64) {

}
