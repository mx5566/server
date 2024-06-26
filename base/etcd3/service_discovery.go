package etcd3

import (
	"context"
	"encoding/json"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"github.com/mx5566/server/base/conf"
	"github.com/mx5566/server/base/entity"
	"github.com/mx5566/server/base/rpc3"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var EndPoints []string = []string{"127.0.0.1"}

// Server/Type/ID  Server/2/crc32(127.0.0.1:9090)
// 服务的发现模块
type ServiceDiscovery struct {
	client *clientv3.Client
}

func NewServiceDiscovery(config conf.ServiceEtcd) *ServiceDiscovery {
	s := &ServiceDiscovery{}
	s.Init(config)

	return s
}

func (sd *ServiceDiscovery) Init(config conf.ServiceEtcd) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: config.EndPoints,
	})
	if err != nil {
		logm.PanicfE("连接etcd3服务器失败:%v, 原因: %s\n", config.EndPoints, err.Error())
		return
	}

	sd.client = client

	sd.Start()
}

func (sd *ServiceDiscovery) Start() {
	go sd.WatchServices()
}

// 服务发现
func (sd *ServiceDiscovery) DiscoverServices() error {
	resp, err := sd.client.Get(context.Background(), base.ServiceName, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	//var buf [4096]byte
	//n := runtime.Stack(buf[:], false)
	//fmt.Printf("==> %s", string(buf[:n]))

	for _, kv := range resp.Kvs {
		sd.AddServiceNode(kv.Value)
	}
	return nil
}

// 监听
func (sd *ServiceDiscovery) WatchServices() error {
	watchChan := sd.client.Watch(context.Background(), base.ServiceName, clientv3.WithPrefix(), clientv3.WithPrevKV())
	sd.DiscoverServices()
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			switch event.Type.String() {
			case "PUT":
				sd.AddServiceNode(event.Kv.Value)
			case "DELETE":
				//logm.DebugfE("ServiceDiscovery WatchServices DeleModule: %s %s cv:%d, mv:%d, vv:%d, ls: %d",
				//	string(event.Kv.Key), string(event.Kv.Value),
				//	event.Kv.CreateRevision, event.Kv.ModRevision,
				//	event.Kv.Version, event.Kv.Lease)
				sd.DeleteServiceNode(event.PrevKv.Value)
			}
		}
	}

	return nil
}

func (sd *ServiceDiscovery) AddServiceNode(by []byte) {
	clsterInfo := &rpc3.ClusterInfo{}

	err := json.Unmarshal(by, clsterInfo)
	if err != nil {
		return
	}

	logm.DebugfE("服务发现新的服务器配置:%s", clsterInfo.String())
	entity.GEntityMgr.SendMsg(rpc3.RpcHead{}, "Cluster.AddClusterNode", clsterInfo)
}

func (sd *ServiceDiscovery) DeleteServiceNode(by []byte) {
	clsterInfo := rpc3.ClusterInfo{}

	err := json.Unmarshal(by, &clsterInfo)
	if err != nil {
		logm.ErrorfE("服务发现注册删除服务节点: %v", clsterInfo)
		return
	}

	entity.GEntityMgr.SendMsg(rpc3.RpcHead{}, "Cluster.DelClusterNode", clsterInfo)

}
