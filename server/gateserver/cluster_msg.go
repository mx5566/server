package gateserver

import (
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
			SERVER.SendToClient()
		}

	} else {
		entity.GEntityMgr.Send(*rpcPacket)
	}
}
