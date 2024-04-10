package gateserver

import (
	"context"
	"github.com/mx5566/server/base/entity"
	"sync"
)

type LoginState uint8

const (
	LoginState_AccountLoging LoginState = iota // 账号登录中
	LoginState_AccountLogin                    // 账号登陆了 进入角色选择界面
	LoginState_PlayerLogin                     // 玩家选择角色登录游戏了
)

type Player struct {
	PlayerID     int64      // 玩家的数据库ID
	ConnID       uint32     // 客户端连接的socketid
	GameServerID uint32     // 所在的游戏服务器ID
	AccountID    int64      // 账号ID数据库ID
	State        LoginState // 登录状态
}

// 管理所有登录的玩家
type PlayerMgr struct {
	entity.Entity
	playersAccount map[int64]*Player // key账号ID
	playersId      map[int64]*Player // key是玩家id
	sync.Mutex
}

func (m *PlayerMgr) Init() {
	m.Entity.Init()
	m.Entity.Start()
	entity.GEntityMgr.RegisterEntity(m)
}

func (m *PlayerMgr) GetPlayer(connId uint32) {

}

func (m *PlayerMgr) AcountLogin(ctx context.Context, accountId int64) {
	m.Mutex.Lock()
	player, ok := m.playersAccount[accountId]
	if ok && player.State == LoginState_AccountLoging {
		player.AccountID = accountId
		player.State = LoginState_AccountLogin
	}
	defer m.Mutex.Unlock()
}

func (m *PlayerMgr) PlayerLogin(ctx context.Context, accountId, playerId int64) {

}
