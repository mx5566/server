package etcd3

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/rpc"
	"go.etcd.io/etcd/clientv3"
)

var EndPoints []string = []string{"127.0.0.1"}

const ServiceName = "Server/"

// Server/Type/ID  Server/2/crc32(127.0.0.1:9090)
// 服务的发现模块
type ServiceDiscovery struct {
	client *clientv3.Client
}

func NewServiceDiscovery(endPoints []string) *ServiceDiscovery {
	s := &ServiceDiscovery{}
	s.Init(endPoints)

	return s
}

func (sd *ServiceDiscovery) Init(endPoints []string) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: endPoints,
	})
	if err != nil {
		logm.PanicfE("连接etcd3服务器失败:%v, 原因: %s\n", endPoints, err.Error())
		return
	}

	sd.client = client
	sd.Start()
	sd.DiscoverServices()
}

func (sd *ServiceDiscovery) Start() {
	go sd.WatchServices()
}

// 服务发现
func (sd *ServiceDiscovery) DiscoverServices() error {
	resp, err := sd.client.Get(context.Background(), ServiceName, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, kv := range resp.Kvs {
		sd.AddServiceNode(kv.Value)
	}
	return nil
}

// 监听
func (sd *ServiceDiscovery) WatchServices() error {
	watchChan := sd.client.Watch(context.Background(), ServiceName, clientv3.WithPrefix())
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			switch event.Type.String() {
			case "PUT":
				sd.AddServiceNode(event.Kv.Value)
			case "DELETE":
				sd.DeleteServiceNode(event.Kv.Value)
			}
		}
	}

	return nil
}

func (sd *ServiceDiscovery) AddServiceNode(by []byte) {
	clsterInfo := rpc.ClusterInfo{}

	err := proto.Unmarshal(by, &clsterInfo)
	if err != nil {
		return
	}

}

func (sd *ServiceDiscovery) DeleteServiceNode(by []byte) {
	clsterInfo := rpc.ClusterInfo{}

	err := proto.Unmarshal(by, &clsterInfo)
	if err != nil {
		return
	}
}
