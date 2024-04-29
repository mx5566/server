package sceneserver

import (
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/cluster"
	"github.com/mx5566/server/base/conf"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/network"
	"github.com/mx5566/server/base/rpc3"
	"github.com/mx5566/server/server/sceneserver/player"
)

type Config struct {
	conf.DB          `yaml:"DB"`
	conf.Server      `yaml:"scene"`
	conf.ModuleEtcd  `yaml:"moduleetcd"`
	conf.ServiceEtcd `yaml:"etcd"`
	conf.Nats        `yaml:"nats"`
	conf.ModuleP     `yaml:"module"`
}

type SceneServer struct {
	s      *network.ServerSocket
	config Config
}

var SERVER SceneServer

func (gs *SceneServer) GetServer() *network.ServerSocket {
	return gs.s
}

func (gs *SceneServer) Init() {
	conf.ReadConf("./config.yaml", &gs.config)

	// 日志初始化
	logm.Init("gateserver", map[string]string{"errFile": "gate_server.log", "logFile": "gate_server_error.log"}, "debug")
	s := new(network.ServerSocket)
	s.Init(gs.config.Server.Ip, gs.config.Server.Port)
	s.Start()

	gs.s = s

	cluster.GCluster.InitCluster(&rpc3.ClusterInfo{
		Ip:          gs.config.Server.Ip,
		Port:        uint32(gs.config.Server.Port),
		ServiceType: rpc3.ServiceType_SceneServer,
	}, gs.config.ServiceEtcd, gs.config.Nats)
	cluster.GCluster.BindPacketFunc(entity.GEntityMgr.PacketFunc)

	gs.InitMgr()

	pr := new(base.Pprof)
	pr.Init()
}

func (gs *SceneServer) InitMgr() {
	player.PLAYERMGR.Init()
}

// 可以用IP+PORT 求一个哈希值
func (gs *SceneServer) GetID() uint32 {
	return cluster.GCluster.ClusterInfo.Id()
}
