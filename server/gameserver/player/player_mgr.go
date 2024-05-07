package player

import (
	"context"
	"github.com/mx5566/server/base/cluster"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/rpc3"
	"reflect"
)

var PLAYERMGR PlayerMgr

type PlayerMgr struct {
	entity.Entity
	entity.EntityPool

	mailBoxAgents map[int64]*cluster.MailBoxAgent
}

func (m *PlayerMgr) Init() {
	m.mailBoxAgents = make(map[int64]*cluster.MailBoxAgent)

	m.Entity.Init()
	m.EntityPool.InitPool(reflect.TypeOf(Player{}))
	entity.GEntityMgr.RegisterEntity(m)
	m.Entity.Start()
}

func (m *PlayerMgr) PlayerLoginRequest(ctx context.Context, accountID, playerId int64, gateClusterId uint32) {
	mMailInfo := cluster.GCluster.GetMailBox(playerId)
	if mMailInfo == nil {
		agent := &cluster.MailBoxAgent{}

		mBox := rpc3.MailBox{}
		mBox.ID = playerId
		mBox.ClusterID = cluster.GCluster.Id()
		mBox.MType = rpc3.MailType_Player
		agent.Init(mBox)
		if !agent.RegisterAgent() {
			return
		}

		mMailInfo = &mBox
		m.mailBoxAgents[playerId] = agent
	}

	head := ctx.Value("rpcHead").(rpc3.RpcHead)

	cluster.GCluster.SendMsg(&rpc3.RpcHead{
		DestServerID: gateClusterId,
		ConnID:       head.ConnID,
		SrcServerID:  mMailInfo.GetClusterID(),
		ID:           accountID},
		"gateserver<-PlayerMgr.PlayerLogin", mMailInfo)

	/*// 角色登录 加载角色数据
	// 数据库加载玩家数据
	filter := mongodb.Newfilter().EQ("playerID", playerId)

	pInstance := mongodb.NewMGDB[model.PlayerSimpleInfo]("game", "player_tbl")
	ops := options.FindOneOptions{}
	ops.SetProjection(bson.D{{"_id", 0}})
	player, err := pInstance.FindOne(context.Background(), filter, ops)
	if err != nil {
		logm.ErrorfE("数据库查找角色失败:%s", err.Error())
		return
	}*/

}

func (m *PlayerMgr) PlayerLogin(ctx context.Context, mailBox *rpc3.MailBox) {
	agent, ok := m.mailBoxAgents[mailBox.ID]
	if !ok {
		return
	}

	head := ctx.Value("rpcHead").(rpc3.RpcHead)

	player := &Player{}
	player.agent = agent
	player.playerId = mailBox.GetID()
	player.SetID(mailBox.GetID())
	player.Init()
	m.AddEntity(player)

	delete(m.mailBoxAgents, mailBox.GetID())

	entity.GEntityMgr.SendMsg(rpc3.RpcHead{}, "Player.Login", head.SrcServerID)
}
