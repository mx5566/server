package gateserver

import (
	"context"
	"github.com/mx5566/server/base/cluster"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/rpc3"
	"sync"
)

type Player struct {
	PlayerID      int64  // 玩家的数据库ID
	ConnID        uint32 // 客户端连接的socketid
	GameServerID  uint32 // 所在的游戏服务器ID
	SceneServerID uint32 // 地图服务器ID
	AccountID     int64
}

var PLAYERMGR PlayerMgr

// 管理所有登录的玩家
type PlayerMgr struct {
	entity.Entity

	connIdPlayers  map[uint32]*Player // 连接网络层生成的ID
	PlayerIDConnID map[int64]uint32   // 玩家id与连接的socketid的映射
	sync.Mutex
}

func (m *PlayerMgr) Init() {
	m.Entity.Init()
	m.Entity.Start()
	entity.GEntityMgr.RegisterEntity(m)

	m.connIdPlayers = make(map[uint32]*Player)
	m.PlayerIDConnID = make(map[int64]uint32)
}

func (m *PlayerMgr) GetPlayer(connId uint32) {

}

func (m *PlayerMgr) AccountLogining(ctx context.Context, player *Player) {

}

func (m *PlayerMgr) PlayerLogin(ctx context.Context, mailBox *rpc3.MailBox) {
	head := ctx.Value("rpcHead").(rpc3.RpcHead)

	playerId := mailBox.ID
	m.Mutex.Lock()
	connId, ok := m.PlayerIDConnID[playerId]
	if ok {
		// 玩家有一个旧的连接
		delete(m.PlayerIDConnID, playerId)
		delete(m.connIdPlayers, connId)
		m.Mutex.Unlock()

		// 服务器主动断开连接
		SERVER.GetServer().StopOneClient(connId)
	}

	p := new(Player)
	p.PlayerID = playerId
	p.AccountID = head.ID
	p.GameServerID = mailBox.ClusterID

	m.Mutex.Lock()
	m.PlayerIDConnID[playerId] = head.ConnID
	m.connIdPlayers[connId] = p
	m.Mutex.Unlock()

	cluster.GCluster.SendMsg(&rpc3.RpcHead{
		DestServerID: head.SrcServerID,
		ConnID:       head.ConnID,
		SrcServerID:  SERVER.GetID(),
		ID:           head.ID},
		"gameserver<-PlayerMgr.PlayerLogin", mailBox)

}

func (m *PlayerMgr) DeletePlayer(ctx context.Context, connId uint32, playerId int64) {
	m.Mutex.Lock()
	delete(m.PlayerIDConnID, playerId)
	delete(m.connIdPlayers, connId)
	m.Mutex.Unlock()
}

func (m *PlayerMgr) GetAccountID(connId uint32) int64 {
	m.Mutex.Lock()
	p, ok := m.connIdPlayers[connId]
	if ok {
		return p.AccountID
	}
	m.Mutex.Unlock()

	return -1
}
