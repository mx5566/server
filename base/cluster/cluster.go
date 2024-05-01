package cluster

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/conf"
	"github.com/mx5566/server/base/consistency"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/etcd3"
	"github.com/mx5566/server/base/network"
	"github.com/mx5566/server/base/rpc3"
	"github.com/mx5566/server/server/pb"
	"github.com/nats-io/nats.go"
	"strings"
	"sync"
)

var GCluster Cluster

const ServerTypeMax = int(rpc3.ServiceType_SceneServer) + 1

type OP struct {
	Module     conf.ModuleP
	ModuleEtcd conf.ModuleEtcd
}

type OpOption func(op *OP)

func (o *OP) Apply(option ...OpOption) {
	for _, v := range option {
		v(o)
	}
}

func WithModuleEtcd(etcd conf.ModuleEtcd, mo conf.ModuleP) OpOption {
	return func(op *OP) {
		op.ModuleEtcd = etcd
		op.Module = mo
	}
}

type Cluster struct {
	entity.Entity
	*rpc3.ClusterInfo
	clusterMap   [ServerTypeMax]map[uint32]*rpc3.ClusterInfo
	clusterMutex sync.Mutex
	hashRingMap  [ServerTypeMax]*consistency.HashRing
	hashMutex    sync.Mutex

	serviceRegister  *etcd3.ServiceRegister
	serviceDiscovery *etcd3.ServiceDiscovery

	// nats client
	natsClient *nats.Conn
	handleFunc network.HandleFunc

	moduleMgr ModuleMgr
	moduleP   conf.ModuleP
}

func (c *Cluster) SendMsg(head *rpc3.RpcHead, funcName string, param ...interface{}) {
	rpcPacket := pb.Marshal(head, &funcName, param...)

	c.Send(rpcPacket)
}

func (c *Cluster) Send(packet rpc3.RpcPacket) {
	head := packet.Head
	switch head.MsgSendType {
	case rpc3.SendType_SendType_Local:
		entity.GEntityMgr.Send(packet)
	case rpc3.SendType_SendType_Single:
		if head.DestServerType == rpc3.ServiceType_WorldServer {
			// 不知道那个服务器
			if head.DestServerID == 0 {
				count := c.moduleP.ModuleCount[head.ClassName]

				index := head.ID % count

				mInfo := c.moduleMgr.GetModule(rpc3.ModuleType(rpc3.ModuleType_value[head.ClassName]), index)
				head.DestServerID = mInfo.ClusterID
			}

			top := fmt.Sprintf("%s%s/%d", base.ServiceName, head.DestServerType.String(), head.DestServerID)

			buff, _ := proto.Marshal(&packet)
			_ = c.natsClient.Publish(top, buff)

			logm.ErrorfE("发送数据到worlsserv: %s", top)
		} else if head.DestServerType == rpc3.ServiceType_GateServer {
			top := fmt.Sprintf("%s%s/%d", base.ServiceName, head.DestServerType.String(), head.DestServerID)

			buff, _ := proto.Marshal(&packet)
			_ = c.natsClient.Publish(top, buff)
		} else if head.DestServerType == rpc3.ServiceType_GameServer {
			top := fmt.Sprintf("%s%s/%d", base.ServiceName, head.DestServerType.String(), head.DestServerID)

			buff, _ := proto.Marshal(&packet)
			_ = c.natsClient.Publish(top, buff)
		}
	case rpc3.SendType_SendType_BroadCast:
		buff, _ := proto.Marshal(&packet)
		c.natsClient.Publish(fmt.Sprintf("%s%s", base.ServiceName, head.DestServerType.String()), buff)
	}
}

func (c *Cluster) RandomClusterByType(serviceType rpc3.ServiceType, id int64) uint32 {
	c.hashMutex.Lock()
	node, _ := c.hashRingMap[serviceType].GetNodeint(id)
	c.hashMutex.Unlock()

	return node

}

func (c *Cluster) NatsSubscibe() {
	if c.natsClient != nil {
		top := c.ClusterInfo.GetTopicServerID()
		logm.DebugfE("订阅的主题1: %s", top)
		c.natsClient.Subscribe(top, func(msg *nats.Msg) {
			packet := rpc3.Packet{
				Id:   0,
				Buff: msg.Data,
			}

			c.HandlePacket(packet)
		})

		top = c.ClusterInfo.GetTopicServerType()
		logm.DebugfE("订阅的主题2: %s", top)

		c.natsClient.Subscribe(top, func(msg *nats.Msg) {
			packet := rpc3.Packet{
				Id:   0,
				Buff: msg.Data,
			}

			c.HandlePacket(packet)
		})
	}
}

func (c *Cluster) InitCluster(clusterInfo *rpc3.ClusterInfo, config conf.ServiceEtcd, natsConfig conf.Nats, params ...OpOption) {
	c.ClusterInfo = clusterInfo

	for i := 0; i < ServerTypeMax; i++ {
		c.clusterMap[i] = make(map[uint32]*rpc3.ClusterInfo)
		c.hashRingMap[i] = consistency.NewRing(nil)
	}

	c.Entity.Init()
	c.Entity.Start()

	entity.RegisterEntity(c)

	op := OP{}
	op.Apply(params...)

	// 服务的注册
	c.serviceRegister = etcd3.NewServiceRegister(clusterInfo, config)
	//服务的发现
	c.serviceDiscovery = etcd3.NewServiceDiscovery(config)

	c.InitNats(natsConfig)

	if len(op.ModuleEtcd.EndPoints) > 0 {
		c.moduleMgr.Init(op.ModuleEtcd.EndPoints, op.ModuleEtcd.GrantTime)
		c.moduleP = op.Module
	}

}

func (c *Cluster) BindPacketFunc(hFunc network.HandleFunc) {
	c.handleFunc = hFunc
}

func (c *Cluster) HandlePacket(packet rpc3.Packet) {
	c.handleFunc(packet)
}

func (c *Cluster) InitNats(natsConfig conf.Nats) bool {
	ops := []nats.Option{}

	op := nats.DisconnectErrHandler(func(conn *nats.Conn, err error) {
		logm.ErrorfE("nats disconnect err:%s", err.Error())
	})
	ops = append(ops, op)

	op = nats.ClosedHandler(func(conn *nats.Conn) {
		logm.ErrorfE("nats close")

	})
	ops = append(ops, op)

	op = nats.MaxReconnects(600)

	ops = append(ops, op)

	op = nats.ReconnectHandler(func(conn *nats.Conn) {
		logm.ErrorfE("nats reconnect")

		c.NatsSubscibe()
	})
	ops = append(ops, op)

	op = nats.ConnectHandler(func(conn *nats.Conn) {
		logm.DebugfE("nats conn")

		c.NatsSubscibe()

		logm.InfofE("连接nats服务器成功url: %s", conn.ConnectedUrl())
	})
	ops = append(ops, op)

	op = nats.RetryOnFailedConnect(true)
	ops = append(ops, op)

	url := strings.Join(natsConfig.EndPoints, ",")

	connect, err := nats.Connect(url, ops...)
	if err != nil {
		logm.ErrorfE("连接nats服务器失败 err : %s", err.Error())
		return false
	}

	c.natsClient = connect

	return true
}

func (c *Cluster) AddClusterNode(ctx context.Context, info *rpc3.ClusterInfo) {
	c.clusterMutex.Lock()
	c.clusterMap[info.GetServiceType()][info.Id()] = info
	c.clusterMutex.Unlock()

	c.hashMutex.Lock()
	c.hashRingMap[info.GetServiceType()].Add(info.Ips())
	c.hashMutex.Unlock()

	logm.InfofE("增加集群信息: %v", info)
}

func (c *Cluster) DelClusterNode(ctx context.Context, info *rpc3.ClusterInfo) {
	c.clusterMutex.Lock()
	delete(c.clusterMap[info.GetServiceType()], info.Id())
	c.clusterMutex.Unlock()

	c.hashMutex.Lock()
	c.hashRingMap[info.GetServiceType()].Remove(info.Ips())
	c.hashMutex.Unlock()

	logm.InfofE("移除集群信息: %v", info)

}

func (c *Cluster) IsEnough(t rpc3.ModuleType) bool {
	n := c.GetModuleMax(t)
	c1 := c.moduleMgr.GetModuleNum(t)
	return c1 >= int(n)

}

func (c *Cluster) GetModuleMax(t rpc3.ModuleType) int64 {
	if _, ok := c.moduleP.ModuleCount[t.String()]; ok {
		return c.moduleP.ModuleCount[t.String()]
	}

	return 0
}
