package gameserver

import (
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base/cluster"
	"github.com/mx5566/server/base/network"
	"github.com/mx5566/server/base/rpc3"
	"github.com/mx5566/server/server/table"
)

type GameServer struct {
}

var SERVER GameServer

func (gs *GameServer) Init() {
	// 日志初始化
	logm.Init("gameserver", map[string]string{"errFile": "game_server.log", "logFile": "game_server_error.log"}, "debug")
	gs.TestLoadTable()

	s := new(network.ServerSocket)
	s.Init("0.0.0.0", 9090)
	s.Start()

	cluster.GCluster.InitCluster(&rpc3.ClusterInfo{
		Ip:          "0.0.0.0",
		Port:        9090,
		ServiceType: rpc3.ServiceType_GameServer,
	}, rpc3.EtcdConfig{
		EndPoints: []string{"127.0.0.1:2379"},
		TimeNum:   10,
	}, rpc3.NatsConfig{
		EndPoints: []string{"127.0.0.1:4222"},
	})
}

func (gs *GameServer) TestLoadTable() {
	table.LoadItemTable("./table/item.xlsx")
}
