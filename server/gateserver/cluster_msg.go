package gateserver

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/mx5566/server/base/cluster"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/rpc3"
)

func HandleMsg(packet rpc3.Packet) {
	rpcPacket := &rpc3.RpcPacket{}
	_ = proto.Unmarshal(packet.Buff, rpcPacket)
	if cluster.GCluster.ClusterInfo.GetServiceType() == rpc3.ServiceType_GateServer {
		// 一种格式需要本地处理 一种是转发到客户端
		if entity.GEntityMgr.IsHasMethod(rpcPacket.Head.ClassName, rpcPacket.Head.FuncName) {
			// 本地有映射的方法
			entity.GEntityMgr.Send(*rpcPacket)
		} else {
			// 需要转发到客户端
			SERVER.SendToClient(*rpcPacket)
		}
	} else {
		entity.GEntityMgr.Send(*rpcPacket)
	}
}

type ClusterMsg struct {
	entity.Entity
}

func (c *ClusterMsg) Init() {
	c.Entity.Init()
	c.Entity.Start()
	entity.GEntityMgr.RegisterEntity(c)
}

// 账号登录了
func (c *ClusterMsg) AccountLogin(ctx context.Context, accountId int64, accountName string) {
	head := ctx.Value("rpcHead").(rpc3.RpcHead)
	entity.GEntityMgr.SendMsg(head, "PlayerMgr.AccountLogin", accountId)
}

func (c ClusterMsg) PlayerLogin(ctx context.Context, accountId, playerId int64) {
	head := ctx.Value("rpcHead").(rpc3.RpcHead)
	entity.GEntityMgr.SendMsg(head, "PlayerMgr.PlayerLogin", accountId, playerId)
}
