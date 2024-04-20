package gateserver

import (
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/cluster"
	"github.com/mx5566/server/base/conf"
	"github.com/mx5566/server/base/network"
	"github.com/mx5566/server/base/rpc3"
)

type Config struct {
	conf.DB          `yaml:"DB"`
	conf.Server      `yaml:"gate"`
	conf.ModuleEtcd  `yaml:"moduleetcd"`
	conf.ServiceEtcd `yaml:"etcd"`
	conf.Nats        `yaml:"nats"`
	conf.ModuleP     `yaml:"module"`
}

type GateServer struct {
	s      *network.ServerSocket
	config Config
}

var SERVER GateServer

func (gs *GateServer) GetServer() *network.ServerSocket {
	return gs.s
}

func (gs *GateServer) Init() {
	conf.ReadConf("./config.yaml", &gs.config)

	// 日志初始化
	logm.Init("gateserver", map[string]string{"errFile": "gate_server.log", "logFile": "gate_server_error.log"}, "debug")
	s := new(network.ServerSocket)
	s.Init(gs.config.Server.Ip, gs.config.Server.Port)

	session := new(ClientSession)
	session.Init()

	s.BindPacketFunc(session.HandlePacket)
	s.Start()

	gs.s = s

	cluster.GCluster.InitCluster(&rpc3.ClusterInfo{
		Ip:          gs.config.Server.Ip,
		Port:        uint32(gs.config.Server.Port),
		ServiceType: rpc3.ServiceType_GateServer,
	}, gs.config.ServiceEtcd, gs.config.Nats)

	cluster.GCluster.BindPacketFunc(HandleMsg)

	gs.InitMgr()

	n := new(ClusterMsg)
	n.Init()

	pr := new(base.Pprof)
	pr.Init()

}

func (gs *GateServer) InitMgr() {
	// 初始化playerMGr
	PLAYERMGR.Init()
}

// 可以用IP+PORT 求一个哈希值
func (gs *GateServer) GetID() uint32 {
	return cluster.GCluster.ClusterInfo.Id()
}

func (gs *GateServer) SendToClient(rpcPacket rpc3.RpcPacket) {

	gs.s.SendMsg(rpcPacket)
}
