package gateserver

import (
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base/cluster"
	"github.com/mx5566/server/base/network"
	"github.com/mx5566/server/base/rpc3"
)

type GateServer struct {
	s *network.ServerSocket
}

var SERVER GateServer

func (gs *GateServer) GetServer() *network.ServerSocket {
	return gs.s
}

func (gs *GateServer) Init() {
	// 日志初始化
	logm.Init("gateserver", map[string]string{"errFile": "gate_server.log", "logFile": "gate_server_error.log"}, "debug")
	s := new(network.ServerSocket)
	s.Init("0.0.0.0", 8080)

	session := new(ClientSession)
	session.Init()

	s.BindPacketFunc(session.HandlePacket)
	s.Start()

	gs.s = s

	cluster.GCluster.InitCluster(&rpc3.ClusterInfo{
		Ip:          "0.0.0.0",
		Port:        8080,
		ServiceType: rpc3.ServiceType_GateServer,
	}, rpc3.EtcdConfig{
		EndPoints: []string{"127.0.0.1:2379"},
		TimeNum:   10,
	}, rpc3.NatsConfig{
		EndPoints: []string{"127.0.0.1:4222"},
	})
	cluster.GCluster.BindPacketFunc(HandleMsg)

}

// 可以用IP+PORT 求一个哈希值
func (gs *GateServer) GetID() uint32 {
	return cluster.GCluster.ClusterInfo.Id()
}

func (gs *GateServer) SendToClient(rpcPacket rpc3.RpcPacket) {

	gs.s.SendMsg(rpcPacket)
}
