package player

import (
	"context"
	"github.com/mx5566/server/base/entity"
	"reflect"
)

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

func (m *PlayerMgr) PlayerLoginRequest(ctx context.Context, playerId int64, gateClusterId uint32, connID uint32) {
	
}
