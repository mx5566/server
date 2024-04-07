package cluster

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/etcd3"
	"github.com/mx5566/server/base/network"
	"github.com/mx5566/server/base/rpc3"
	"github.com/mx5566/server/server/pb"
	"github.com/nats-io/nats.go"
	"hash/crc32"
	"strings"
	"sync"
)

var GCluster Cluster

const ServerTypeMax = int(rpc3.ServiceType_WorldServer) + 1

type Cluster struct {
	entity.Entity
	*rpc3.ClusterInfo
	clusterMap   [ServerTypeMax]map[uint32]*rpc3.ClusterInfo
	clusterMutex sync.Mutex

	serviceRegister  *etcd3.ServiceRegister
	serviceDiscovery *etcd3.ServiceDiscovery

	natsClient *nats.Conn
	handleFunc network.HandleFunc
}

func (c *Cluster) SendMsg(head *rpc3.RpcHead, funcName string, param ...interface{}) {
	rpcPacket := pb.Marshal(head, &funcName, param)

	switch head.MsgSendType {
	case rpc3.SendType_SendType_Local:
		entity.GEntityMgr.Send(rpcPacket)
	case rpc3.SendType_SendType_Single:
		if head.DestServerType == rpc3.ServiceType_WorldServer {
			// 有多个worldserver 发送那个呢
			// 先随机一个服务器
			clusterID := c.clusterMap[head.DestServerType][crc32.ChecksumIEEE([]byte("0.0.0.0:9999"))]

			top := fmt.Sprintf("%s%s/%d", base.ServiceName, head.DestServerType.String(), clusterID.Id())

			head.DestServerID = clusterID.Id()

			buff, _ := proto.Marshal(&rpcPacket)
			_ = c.natsClient.Publish(top, buff)

			logm.ErrorfE("发送数据到worlsserv: %s", top)
		} else if head.DestServerType == rpc3.ServiceType_GateServer {
			top := fmt.Sprintf("%s%s/%d", base.ServiceName, head.DestServerType.String(), head.DestServerID)

			buff, _ := proto.Marshal(&rpcPacket)
			_ = c.natsClient.Publish(top, buff)
		} else if head.DestServerType == rpc3.ServiceType_GameServer {

		}
	case rpc3.SendType_SendType_BroadCast:
		buff, _ := proto.Marshal(&rpcPacket)
		c.natsClient.Publish(fmt.Sprintf("%s%s", base.ServiceName, head.DestServerType.String()), buff)
	}

}

func (c *Cluster) RandomClusterBuType(serviceType rpc3.ServiceType) *rpc3.ClusterInfo {
	return nil

}

func (c *Cluster) InitCluster(clusterInfo *rpc3.ClusterInfo, config rpc3.EtcdConfig, natsConfig rpc3.NatsConfig) {
	c.ClusterInfo = clusterInfo

	for i := 0; i < ServerTypeMax; i++ {
		c.clusterMap[i] = make(map[uint32]*rpc3.ClusterInfo)
	}

	c.Entity.Init()
	c.Entity.Start()

	entity.RegisterEntity(c)

	// 服务的注册
	c.serviceRegister = etcd3.NewServiceRegister(clusterInfo, config)
	//服务的发现
	c.serviceDiscovery = etcd3.NewServiceDiscovery(config)

	c.InitNats(natsConfig)

	if c.natsClient != nil {
		top := fmt.Sprintf("%s%s/%d", base.ServiceName, c.ClusterInfo.GetServiceType().String(), c.ClusterInfo.Id())
		logm.DebugfE("订阅的主题1: %s", top)
		c.natsClient.Subscribe(top, func(msg *nats.Msg) {
			packet := rpc3.Packet{
				Id:   0,
				Buff: msg.Data,
			}

			c.HandlePacket(packet)
		})

		top = fmt.Sprintf("%s%s", base.ServiceName, c.ClusterInfo.GetServiceType().String())
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

func (c *Cluster) BindPacketFunc(hFunc network.HandleFunc) {
	c.handleFunc = hFunc
}

func (c *Cluster) HandlePacket(packet rpc3.Packet) {
	c.handleFunc(packet)
}

func (c *Cluster) HandleMsg(packet rpc3.Packet) {
	rpcPacket := &rpc3.RpcPacket{}
	_ = proto.Unmarshal(packet.Buff, rpcPacket)
	if c.ClusterInfo.GetServiceType() == rpc3.ServiceType_GateServer {
		// 一种格式需要本地处理 一种是转发到客户端
		if entity.GEntityMgr.IsHasMethod(rpcPacket.Head.ClassName, rpcPacket.Head.FuncName) {
			// 本地有映射的方法
			entity.GEntityMgr.Send(*rpcPacket)
		} else {
			// 需要转发到客户端

		}

	} else {
		entity.GEntityMgr.Send(*rpcPacket)
	}
}

func (c *Cluster) InitNats(natsConfig rpc3.NatsConfig) {
	ops := []nats.Option{}

	op := nats.DisconnectErrHandler(func(conn *nats.Conn, err error) {
		logm.ErrorfE("nats disconnect err:%s", err.Error())
	})
	ops = append(ops, op)

	op = nats.ClosedHandler(func(conn *nats.Conn) {
		logm.ErrorfE("nats close")

	})
	ops = append(ops, op)

	op = nats.ReconnectHandler(func(conn *nats.Conn) {
		logm.ErrorfE("nats reconnect")

	})
	ops = append(ops, op)

	url := strings.Join(natsConfig.GetEndPoints(), ",")

	connect, err := nats.Connect(url, ops...)
	if err != nil {
		logm.ErrorfE("连接nats服务器失败 err : %s", err.Error())
		return
	}

	c.natsClient = connect
	logm.InfofE("连接nats服务器成功url: %s", url)
}

func (c *Cluster) AddClusterNode(ctx context.Context, info *rpc3.ClusterInfo) {
	c.clusterMutex.Lock()
	c.clusterMap[info.GetServiceType()][info.Id()] = info
	c.clusterMutex.Unlock()

	logm.InfofE("增加集群信息: %v", c.clusterMap)
}

func (c *Cluster) DelClusterNode(ctx context.Context, info *rpc3.ClusterInfo) {
	c.clusterMutex.Lock()
	delete(c.clusterMap[info.GetServiceType()], info.Id())
	c.clusterMutex.Unlock()
}
