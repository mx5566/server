package gameserver

import (
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base/cluster"
	"github.com/mx5566/server/base/conf"
	"github.com/mx5566/server/base/network"
	"github.com/mx5566/server/base/rpc3"
	"github.com/mx5566/server/server/table"
)

type Config struct {
	conf.DB          `yaml:"DB"`
	conf.Server      `yaml:"game"`
	conf.ModuleEtcd  `yaml:"moduleetcd"`
	conf.ServiceEtcd `yaml:"etcd"`
	conf.Nats        `yaml:"nats"`
	conf.ModuleP     `yaml:"module"`
}

type GameServer struct {
	s      *network.ServerSocket
	config Config
}

var SERVER GameServer

func (gs *GameServer) Init() {
	conf.ReadConf("./config.yaml", &gs.config)

	// 日志初始化
	logm.Init("gameserver", map[string]string{"errFile": "game_server.log", "logFile": "game_server_error.log"}, "debug")
	gs.TestLoadTable()

	s := new(network.ServerSocket)
	s.Init(gs.config.Server.Ip, gs.config.Server.Port)
	s.Start()

	cluster.GCluster.InitCluster(&rpc3.ClusterInfo{
		Ip:          gs.config.Server.Ip,
		Port:        uint32(gs.config.Server.Port),
		ServiceType: rpc3.ServiceType_GameServer,
	}, gs.config.ServiceEtcd, gs.config.Nats, cluster.WithModuleEtcd(gs.config.ModuleEtcd, gs.config.ModuleP))

	gs.InitMgr()
}

func (gs *GameServer) InitMgr() {

}

func (gs *GameServer) TestLoadTable() {
	table.LoadItemTable("./table/item.xlsx")
}
