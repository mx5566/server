package worldserver

import (
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base/cluster"
	"github.com/mx5566/server/base/conf"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/network"
	"github.com/mx5566/server/base/orm"
	"github.com/mx5566/server/base/rpc3"
	"github.com/mx5566/server/server/worldserver/account"
)

type Config struct {
	conf.Server      `yaml:"world"`
	conf.ModuleEtcd  `yaml:"moduleetcd"`
	conf.ServiceEtcd `yaml:"etcd"`
	conf.Nats        `yaml:"nats"`
	conf.ModuleP     `yaml:"module"`
	conf.DB          `yaml:"DB"`
	conf.MailBoxEtcd `yaml:"mailboxetcd"`
}

type WorldServer struct {
	s      *network.ServerSocket
	config Config
}

var SERVER WorldServer

func (gs *WorldServer) Init() {
	// 日志初始化
	logm.Init("worldserver", map[string]string{"errFile": "world_server_error.log", "logFile": "world_server.log"}, "debug")

	// 配置文件加载
	conf.ReadConf("./config.yaml", &gs.config)

	// 数据库
	orm.OpenMongodb(gs.config.DB)

	s := new(network.ServerSocket)
	s.Init(gs.config.Server.Ip, gs.config.Server.Port)
	s.Start()

	gs.s = s

	cluster.GCluster.InitCluster(&rpc3.ClusterInfo{
		Ip:          gs.config.Server.Ip,
		Port:        uint32(gs.config.Server.Port),
		ServiceType: rpc3.ServiceType_WorldServer,
	}, gs.config.ServiceEtcd,
		gs.config.Nats,
		cluster.WithModuleEtcd(gs.config.ModuleEtcd, gs.config.ModuleP),
		cluster.WithMailBoxEtcd(gs.config.MailBoxEtcd))

	cluster.GCluster.BindPacketFunc(entity.GEntityMgr.PacketFunc)

	// 初始化逻辑
	gs.InitMgr()
}

func (gs *WorldServer) InitMgr() {
	account.GAccountMgr.Init()
}

// 可以用IP+PORT 求一个哈希值
func (gs *WorldServer) GetID() uint32 {
	return cluster.GCluster.ClusterInfo.Id()
}
