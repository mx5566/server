package worldserver

import (
	"github.com/mx5566/logm"
	"github.com/mx5566/server/cluster"
	"github.com/mx5566/server/entity"
	"github.com/mx5566/server/network"
	"github.com/mx5566/server/rpc3"
)

type WorldServer struct {
}

var SERVER WorldServer

func (gs *WorldServer) Init() {
	// 日志初始化
	logm.Init("worldserver", map[string]string{"errFile": "world_server.log", "logFile": "world_server_error.log"}, "debug")
	s := new(network.ServerSocket)
	s.Init("0.0.0.0", 9999)
	s.Start()

	cluster.GCluster.InitCluster(&rpc3.ClusterInfo{
		Ip:          "0.0.0.0",
		Port:        9999,
		ServiceType: rpc3.ServiceType_WorldServer,
	}, rpc3.EtcdConfig{
		EndPoints: []string{"127.0.0.1:2379"},
		TimeNum:   10,
	}, rpc3.NatsConfig{
		EndPoints: []string{"127.0.0.1:4222"},
	})

	cluster.GCluster.BindPacketFunc(entity.GEntityMgr.PacketFunc)
}
