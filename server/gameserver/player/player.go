package player

import (
	"context"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base/cluster"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/orm/mongodb"
	"github.com/mx5566/server/base/rpc3"
	"github.com/mx5566/server/server/gameserver"
	"github.com/mx5566/server/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Player struct {
	entity.Entity
	agent         *cluster.MailBoxAgent
	online        bool // 是否在线
	GateServerId  uint32
	SceneServerId uint32
	playerId      int64
	playerData    model.PlayerData
}

func (p *Player) Init() {
	p.Entity.Init()
	p.Entity.Start()
}

func (p *Player) updateLease() {
	if p.online {
		p.agent.Lease()
	} else {
		p.agent.Delete()
	}
}

func (p *Player) Login(gateserverID uint32) {
	p.GateServerId = gateserverID

	// 数据库加载
	filter := mongodb.Newfilter().EQ("playerID", p.playerId)
	pInstance := mongodb.NewMGDB[model.PlayerData]("game", "player_tbl")
	ops := options.FindOneOptions{}
	ops.SetProjection(bson.D{{"_id", 0}})
	player, err := pInstance.FindOne(context.Background(), filter, &ops)
	if err != nil {
		logm.ErrorfE("数据库查找角色失败:%s", err.Error())
		return
	}

	p.playerData = player

	cluster.GCluster.SendMsg(&rpc3.RpcHead{
		SrcServerID:  gameserver.SERVER.GetID(),
		DestServerID: cluster.GCluster.RandomClusterByType(rpc3.ServiceType_SceneServer, p.playerId),
		ID:           p.playerId,
	}, "sceneserver<-MapMgr.EnterMap",
		999, p.playerId, p.GateServerId, gameserver.SERVER.GetID())
	// 进入地图
}
