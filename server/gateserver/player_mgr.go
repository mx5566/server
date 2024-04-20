package gateserver

import (
	"context"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/rpc3"
	"sync"
)

type Player struct {
	PlayerID     int64           // 玩家的数据库ID
	ConnID       uint32          // 客户端连接的socketid
	GameServerID uint32          // 所在的游戏服务器ID
	AccountID    int64           // 账号ID数据库ID
	State        base.LoginState // 登录状态
}

var PLAYERMGR PlayerMgr

// 管理所有登录的玩家
type PlayerMgr struct {
	entity.Entity

	playerConnId map[uint32]*Player // 连接网络层生成的ID
	sync.Mutex
}

func (m *PlayerMgr) Init() {
	m.Entity.Init()
	m.Entity.Start()
	entity.GEntityMgr.RegisterEntity(m)

	m.playerConnId = make(map[uint32]*Player)
}

func (m *PlayerMgr) GetPlayer(connId uint32) {

}

func (m *PlayerMgr) GetLoginState(connId uint32) base.LoginState {
	player, ok := m.playerConnId[connId]
	if ok && player != nil {
		return player.State
	}

	return base.LoginState_None
}

func (m *PlayerMgr) AccountLogining(ctx context.Context, player *Player) {
	head := ctx.Value("rpcHead").(rpc3.RpcHead)

	m.Mutex.Lock()
	_, ok := m.playerConnId[head.ConnID]
	if !ok {
		m.playerConnId[player.ConnID] = player
	}

	m.Mutex.Unlock()
}

func (m *PlayerMgr) AcountLogin(ctx context.Context, accountId int64) {
	head := ctx.Value("rpcHead").(rpc3.RpcHead)

	m.Mutex.Lock()
	player, ok := m.playerConnId[head.ConnID]
	if ok && player.State == base.LoginState_AccountLogining {
		player.AccountID = accountId
		player.State = base.LoginState_AccountLogin
	}
	defer m.Mutex.Unlock()
}

func (m *PlayerMgr) PlayerLogin(ctx context.Context, accountId, playerId int64) {
	head := ctx.Value("rpcHead").(rpc3.RpcHead)

	m.Mutex.Lock()
	player, ok := m.playerConnId[head.ConnID]
	if ok && player.State == base.LoginState_AccountLogin {
		player.PlayerID = playerId
		player.State = base.LoginState_PlayerLogin
	}

	m.Mutex.Unlock()
}

func (m *PlayerMgr) DeletePlayer(ctx context.Context, connId uint32) {
	m.Mutex.Lock()

	delete(m.playerConnId, connId)

	m.Mutex.Unlock()
}
