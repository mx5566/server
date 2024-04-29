package player

import (
	"context"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/rpc3"
	"reflect"
)

var PLAYERMGR PlayerMgr

type PlayerMgr struct {
	entity.Entity
	entity.EntityPool
}

func (m *PlayerMgr) Init() {
	m.Entity.Init()
	m.EntityPool.InitPool(reflect.TypeOf(Player{}))
	entity.GEntityMgr.RegisterEntity(m)
	m.Entity.Start()
}

func (m *PlayerMgr) LoginMap(ctx context.Context, mapID uint32, playerID int64, gateClusterID, gameClusterID uint32) {
	player := m.GetEntity(playerID)
	if player == nil {
		// 玩家不存在 加载玩家
		p := new(Player)
		p.Unit.ID = playerID
		p.I = p
		p.Init()

		m.AddEntity(p)
	}

	entity.GEntityMgr.SendMsg(rpc3.RpcHead{ID: playerID}, "Map.EnterMap", mapID, playerID, gateClusterID, gameClusterID)
}
